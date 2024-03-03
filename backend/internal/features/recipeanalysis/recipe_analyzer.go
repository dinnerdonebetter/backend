package recipeanalysis

import (
	"context"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/dustin/go-humanize/english"
	"github.com/goccy/go-graphviz"
	"github.com/hako/durafmt"
	"github.com/heimdalr/dag"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

var (
	_ graph.Node = (*recipeStepGraphNode)(nil)

	errNotAcyclic = errors.New("recipe is not acyclic")
)

// recipeStepGraphNode is a node in an implicit graph.
type recipeStepGraphNode struct {
	id int64
}

// newGraphNode returns a new recipeStepGraphNode.
func newGraphNode(id int64) graph.Node {
	return &recipeStepGraphNode{id: id}
}

func (g recipeStepGraphNode) ID() int64 {
	return g.id
}

var errRecipeStepIDNotFound = errors.New("recipe step ID not found")

func findStepIndexForRecipeStepProductID(recipe *types.Recipe, recipeStepProductID string) (int64, error) {
	if step := recipe.FindStepForRecipeStepProductID(recipeStepProductID); step != nil {
		return graphIDForStep(step), nil
	}

	return -1, errRecipeStepIDNotFound
}

func findStepIDForRecipeStepProductID(recipe *types.Recipe, recipeStepProductID string) (string, error) {
	if step := recipe.FindStepForRecipeStepProductID(recipeStepProductID); step != nil {
		return fmt.Sprintf("%d", graphIDForStep(step)), nil
	}

	return "", errRecipeStepIDNotFound
}

func findStepIndexForRecipeStepID(recipe *types.Recipe, recipeStepID string) (int64, error) {
	for _, step := range recipe.Steps {
		if step.ID == recipeStepID {
			return graphIDForStep(step), nil
		}
	}

	return -1, errRecipeStepIDNotFound
}

func graphIDForStep(step *types.RecipeStep) int64 {
	return int64(step.Index + 1)
}

// RecipeAnalyzer analyzes recipes for insights (ugh).
type RecipeAnalyzer interface {
	MakeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (*simple.DirectedGraph, error)
	GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error)
	GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error)
	RenderMermaidDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string
	RenderGraphvizDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string
}

var _ RecipeAnalyzer = (*recipeAnalyzer)(nil)

// recipeAnalyzer creates graphs from recipes.
type recipeAnalyzer struct {
	tracer tracing.Tracer
	logger logging.Logger
}

// NewRecipeAnalyzer creates a recipeAnalyzer.
func NewRecipeAnalyzer(logger logging.Logger, tracerProvider tracing.TracerProvider) RecipeAnalyzer {
	return &recipeAnalyzer{
		logger: logging.EnsureLogger(logger).WithName("recipe_analyzer"),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("recipe_grapher")),
	}
}

// GenerateDAGDiagramForRecipe generates DAG diagrams for a given recipe.
func (g *recipeAnalyzer) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	img, err := g.renderGraph(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (g *recipeAnalyzer) MakeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (*simple.DirectedGraph, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := simple.NewDirectedGraph()

	for _, step := range recipe.Steps {
		recipeGraph.AddNode(newGraphNode(graphIDForStep(step)))
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		for _, vessel := range step.Vessels {
			if vessel.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *vessel.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}
	}

	directedCycles := topo.DirectedCyclesIn(recipeGraph)
	if len(directedCycles) > 0 {
		return nil, errNotAcyclic
	}

	if _, err := g.makeDAGForRecipe(ctx, recipe); err != nil {
		return nil, fmt.Errorf("parsing recipe as DAG: %w", err)
	}

	return recipeGraph, nil
}

type RecipeStepIdentifier struct {
	recipeStep *types.RecipeStep
}

func (i *RecipeStepIdentifier) ID() string {
	return fmt.Sprintf("%d", graphIDForStep(i.recipeStep))
}

// makeDAGForRecipe makes a proper DAG for the provided Recipe.
func (g *recipeAnalyzer) makeDAGForRecipe(ctx context.Context, recipe *types.Recipe) (*dag.DAG, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := dag.NewDAG()

	for _, step := range recipe.Steps {
		if _, err := recipeGraph.AddVertex(&RecipeStepIdentifier{recipeStep: step}); err != nil {
			return nil, fmt.Errorf("adding initial step %d to graph: %w", step.Index, err)
		}
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID == nil {
				continue
			}

			toStepID, err := findStepIDForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step product ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(fmt.Sprintf("%d", graphIDForStep(step)), toStepID); err != nil {
				return nil, fmt.Errorf("adding recipe step %d to graph: %w", step.Index, err)
			}
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID == nil {
				continue
			}

			toStep, err := findStepIDForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step ingredient ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(fmt.Sprintf("%d", graphIDForStep(step)), toStep); err != nil {
				var dupeErr dag.EdgeDuplicateError
				if errors.As(err, &dupeErr) {
					continue
				}

				return nil, fmt.Errorf("adding ingredient step %d to graph: %w", step.Index, err)
			}
		}
	}

	recipeGraph.ReduceTransitively()

	return recipeGraph, nil
}

func (g *recipeAnalyzer) renderGraph(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	gv := graphviz.New()
	pg, err := graphviz.ParseBytes([]byte(g.RenderGraphvizDiagramForRecipe(ctx, recipe)))
	if err != nil {
		return nil, err
	}

	img, err := gv.RenderImage(pg)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// frozenIngredientDefrostStepsFilter iterates through a recipe and returns
// the list of ingredients within that are indicated as kept frozen.
func frozenIngredientDefrostStepsFilter(recipe *types.Recipe) map[string][]int {
	out := map[string][]int{}

	for _, recipeStep := range recipe.Steps {
		ingredientIndices := []int{}
		for i, ingredient := range recipeStep.Ingredients {
			// if it's a valid ingredient
			if ingredient.Ingredient != nil &&
				// if the ingredient has storage temperature set
				ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius != nil &&
				// the ingredient's storage temperature is set to something about freezing temperature.
				*ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius <= 3 {
				ingredientIndices = append(ingredientIndices, i)
			}
		}

		if len(ingredientIndices) > 0 {
			out[recipeStep.ID] = ingredientIndices
		}
	}

	return out
}

func buildThawStepCreationExplanation(recipeStepIndex int64, ingredientIndices ...int) string {
	if len(ingredientIndices) == 0 {
		return ""
	}

	stringIndices := []string{}
	for _, i := range ingredientIndices {
		stringIndices = append(stringIndices, fmt.Sprintf("#%d", i+1))
	}

	var d string
	if len(ingredientIndices) > 1 {
		d = "ingredients"
	} else if len(ingredientIndices) == 1 {
		d = "ingredient"
	}

	return fmt.Sprintf("frozen %s (%s) for step #%d might need to be thawed ahead of time", d, strings.Join(stringIndices, ", "), recipeStepIndex)
}

func (g *recipeAnalyzer) generateMealPlanTasksForFrozenIngredients(ctx context.Context, mealPlanOptionID string, recipe *types.Recipe) []*types.MealPlanTaskDatabaseCreationInput {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.Clone().WithValue(keys.RecipeIDKey, recipe.ID)

	frozenIngredientSteps := frozenIngredientDefrostStepsFilter(recipe)
	logger.WithValue("frozen_steps_qty", len(frozenIngredientSteps)).Info("creating frozen stepSet inputs")

	outputs := []*types.MealPlanTaskDatabaseCreationInput{}
	for stepID, ingredientIndices := range frozenIngredientSteps {
		stepIndex, err := findStepIndexForRecipeStepID(recipe, stepID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "determining recipe stepSet index for stepSet ID")
			continue
		}

		explanation := buildThawStepCreationExplanation(stepIndex, ingredientIndices...)
		if explanation == "" {
			continue
		}

		outputs = append(outputs, &types.MealPlanTaskDatabaseCreationInput{
			ID:                  identifiers.New(),
			CreationExplanation: explanation,
			MealPlanOptionID:    mealPlanOptionID,
		})
	}

	return outputs
}

const recipeTaskStepCreationExplanation = "recipe prep task exists for steps"

func (g *recipeAnalyzer) GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	inputs := g.generateMealPlanTasksForFrozenIngredients(ctx, mealPlanOptionID, recipe)

	for _, prepTask := range recipe.PrepTasks {
		inputs = append(inputs, &types.MealPlanTaskDatabaseCreationInput{
			AssignedToUser:      nil,
			CreationExplanation: recipeTaskStepCreationExplanation,
			StatusExplanation:   "",
			MealPlanOptionID:    mealPlanOptionID,
			RecipePrepTaskID:    prepTask.ID,
			ID:                  identifiers.New(),
		})
	}

	return inputs, nil
}

type provisionCount struct {
	ingredients, instruments, vessels uint
}

func stepProvidesWhatToOtherStep(recipe *types.Recipe, fromStepIndex, toStepIndex uint) string {
	from, to := recipe.Steps[fromStepIndex], recipe.Steps[toStepIndex]
	provides := []string{}

	count := provisionCount{}
	for _, product := range from.Products {
		for _, step := range recipe.Steps {
			if step.ID != to.ID {
				continue
			}

			for _, ingredient := range step.Ingredients {
				if ingredient.RecipeStepProductID != nil && *ingredient.RecipeStepProductID == product.ID {
					count.ingredients++
				}
			}

			for _, instrument := range step.Instruments {
				if instrument.RecipeStepProductID != nil && *instrument.RecipeStepProductID == product.ID {
					count.instruments++
				}
			}

			for _, vessel := range step.Vessels {
				if vessel.RecipeStepProductID != nil && *vessel.RecipeStepProductID == product.ID {
					count.vessels++
				}
			}
		}
	}

	renderCount := func(x uint, typ string) string {
		/*
			unnecessary Sprintf, but I might do something like this later:

			var prefix string
			if x == 1 {
				prefix = "an"
			}
		*/

		return strings.TrimSpace(fmt.Sprintf(" %s ", english.PluralWord(int(x), typ, fmt.Sprintf("%ss", typ))))
	}

	if count.ingredients > 0 {
		provides = append(provides, renderCount(count.ingredients, "ingredient"))
	}

	if count.instruments > 0 {
		provides = append(provides, renderCount(count.ingredients, "instrument"))
	}

	if count.vessels > 0 {
		provides = append(provides, renderCount(count.vessels, "vessel"))
	}

	return english.OxfordWordSeries(provides, "and")
}

func (g *recipeAnalyzer) RenderMermaidDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	var mermaid strings.Builder
	mermaid.WriteString("flowchart TD;\n")

	for _, step := range recipe.Steps {
		mermaid.WriteString(fmt.Sprintf("	Step%d[\"Step #%d (%s)\"];\n", graphIDForStep(step), graphIDForStep(step), step.Preparation.Name))
	}

	for i := range recipe.Steps {
		for j := range recipe.Steps {
			if i == j {
				continue
			}

			if provides := stepProvidesWhatToOtherStep(recipe, uint(i), uint(j)); provides != "" {
				mermaid.WriteString(fmt.Sprintf("\tStep%d -->|%s| Step%d;\n", graphIDForStep(recipe.Steps[i]), provides, graphIDForStep(recipe.Steps[j])))
			}
		}
	}

	for i := range recipe.PrepTasks {
		prepTask := recipe.PrepTasks[i]

		mermaid.WriteString(fmt.Sprintf("subgraph %d [\"%s (prep task #%d)\"]\n", i, prepTask.Name, i+1))
		for j := range prepTask.TaskSteps {
			mermaid.WriteString(fmt.Sprintf("Step%d\n", recipe.FindStepIndexByID(prepTask.TaskSteps[j].BelongsToRecipeStep)))
		}
		mermaid.WriteString("end\n")
	}

	return mermaid.String()
}

var colorNames = []string{
	"CornflowerBlue",
	"purple",
	"Salmon",
	"PeachPuff",
	"Tomato",
	"Orange",
	"Chocolate",
	"LemonChiffon",
}

const graphvizPreamble = `strict digraph {
	fontname="Outfit,Helvetica,Arial,sans-serif";
	node [fontname="Helvetica,Arial,sans-serif" shape=egg];
	edge [fontname="Helvetica,Arial,sans-serif"];
    concentrate=true;
    compound=true;
	rankdir=TB;
	start=1;
`

func (g *recipeAnalyzer) RenderGraphvizDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	var graphViz strings.Builder
	graphViz.WriteString(graphvizPreamble + "\n")

	for _, step := range recipe.Steps {
		graphViz.WriteString(fmt.Sprintf("\tStep%d [label=\"Step #%d (%s)\"];\n", graphIDForStep(step), graphIDForStep(step), step.Preparation.Name))
	}

	for i := range recipe.Steps {
		for j := range recipe.Steps {
			if i == j {
				continue
			}

			if provides := stepProvidesWhatToOtherStep(recipe, uint(i), uint(j)); provides != "" {
				stepLabel := ""
				if recipe.Steps[i].MinimumEstimatedTimeInSeconds != nil && *recipe.Steps[i].MinimumEstimatedTimeInSeconds > 0 {
					stepLabel = durafmt.Parse(time.Duration(*recipe.Steps[i].MinimumEstimatedTimeInSeconds) * time.Second).String()
				}

				graphViz.WriteString(fmt.Sprintf("\tStep%d -> Step%d [color=\"black\" label=%q];\n", graphIDForStep(recipe.Steps[i]), graphIDForStep(recipe.Steps[j]), stepLabel))
			}
		}
	}

	for i := range recipe.PrepTasks {
		prepTask := recipe.PrepTasks[i]

		graphViz.WriteString(fmt.Sprintf("\n\tsubgraph cluster_%d {\n\t\tnode [style=filled];\n\t\t", i))
		for j := range prepTask.TaskSteps {
			var arrow string
			if j != len(prepTask.TaskSteps)-1 {
				arrow = " -> "
			}
			graphViz.WriteString(fmt.Sprintf("Step%d%s", graphIDForStep(recipe.Steps[recipe.FindStepIndexByID(prepTask.TaskSteps[j].BelongsToRecipeStep)]), arrow))
		}
		colorToUse := colorNames[i%len(colorNames)]

		graphViz.WriteString(fmt.Sprintf(";\n\t\tlabel = \"(prep task #%d)\";\n\t\tcolor=%s;\n\t}\n", i+1, colorToUse))
	}

	graphViz.WriteString("}\n")

	return graphViz.String()
}

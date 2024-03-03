package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

func BuildFakeRecipePrepTask() *types.RecipePrepTask {
	recipePrepTaskSteps := []*types.RecipePrepTaskStep{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTaskSteps = append(recipePrepTaskSteps, BuildFakeRecipePrepTaskStep())
	}

	minStorageTemp := buildFakeNumber()
	minTimeBuffer := buildFakeNumber()

	return &types.RecipePrepTask{
		ID:                          BuildFakeID(),
		Notes:                       buildUniqueString(),
		Name:                        buildUniqueString(),
		Description:                 buildUniqueString(),
		Optional:                    fake.Bool(),
		ExplicitStorageInstructions: buildUniqueString(),
		StorageType: fake.RandomString([]string{
			types.RecipePrepTaskStorageTypeUncovered,
			types.RecipePrepTaskStorageTypeCovered,
			types.RecipePrepTaskStorageTypeAirtightContainer,
			types.RecipePrepTaskStorageTypeWireRack,
		}),
		TaskSteps:                              recipePrepTaskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MaximumTimeBufferBeforeRecipeInSeconds: pointer.To(uint32(minTimeBuffer + 1)),
		MinimumStorageTemperatureInCelsius:     pointer.To(float32(minStorageTemp)),
		MaximumStorageTemperatureInCelsius:     pointer.To(float32(minStorageTemp + 1)),
		BelongsToRecipe:                        BuildFakeID(),
		CreatedAt:                              BuildFakeTime(),
		LastUpdatedAt:                          nil,
		ArchivedAt:                             nil,
	}
}

func BuildFakeRecipePrepTaskList() *types.QueryFilteredResult[types.RecipePrepTask] {
	recipePrepTasks := &types.QueryFilteredResult[types.RecipePrepTask]{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTasks.Data = append(recipePrepTasks.Data, BuildFakeRecipePrepTask())
	}

	return recipePrepTasks
}

func BuildFakeRecipePrepTaskStep() *types.RecipePrepTaskStep {
	return &types.RecipePrepTaskStep{
		ID:                      BuildFakeID(),
		BelongsToRecipeStep:     BuildFakeID(),
		BelongsToRecipePrepTask: BuildFakeID(),
		SatisfiesRecipeStep:     fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepCreationRequestInput() *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep: BuildFakeID(),
		SatisfiesRecipeStep: fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepUpdateRequestInput() *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		BelongsToRecipeStep:     pointer.To(BuildFakeID()),
		BelongsToRecipePrepTask: pointer.To(BuildFakeID()),
		SatisfiesRecipeStep:     pointer.To(fake.Bool()),
	}
}

func BuildFakeRecipePrepTaskCreationRequestInput() *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepCreationRequestInput())
	}

	minStorageTemp := buildFakeNumber()
	minTimeBuffer := buildFakeNumber()

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                                  buildUniqueString(),
		ExplicitStorageInstructions:            buildUniqueString(),
		Name:                                   buildUniqueString(),
		Optional:                               fake.Bool(),
		Description:                            buildUniqueString(),
		StorageType:                            types.RecipePrepTaskStorageTypeUncovered,
		BelongsToRecipe:                        BuildFakeID(),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointer.To(uint32(minTimeBuffer) + 1),
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MinimumStorageTemperatureInCelsius:     pointer.To(float32(minStorageTemp)),
		MaximumStorageTemperatureInCelsius:     pointer.To(float32(minStorageTemp + 1)),
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInput() *types.RecipePrepTaskUpdateRequestInput {
	minTemp := buildFakeNumber()

	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInput())
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointer.To(buildUniqueString()),
		ExplicitStorageInstructions:            pointer.To(buildUniqueString()),
		Name:                                   pointer.To(buildUniqueString()),
		Description:                            pointer.To(buildUniqueString()),
		Optional:                               pointer.To(fake.Bool()),
		StorageType:                            pointer.To(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointer.To(BuildFakeID()),
		MaximumTimeBufferBeforeRecipeInSeconds: pointer.To(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointer.To(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointer.To(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointer.To(fake.Uint32()),
		TaskSteps:                              taskSteps,
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInputFromRecipePrepTask(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	minTemp := buildFakeNumber()

	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, converters.ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(x))
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointer.To(buildUniqueString()),
		ExplicitStorageInstructions:            pointer.To(buildUniqueString()),
		Name:                                   pointer.To(buildUniqueString()),
		Description:                            pointer.To(buildUniqueString()),
		Optional:                               pointer.To(fake.Bool()),
		StorageType:                            pointer.To(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointer.To(BuildFakeID()),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointer.To(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointer.To(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointer.To(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointer.To(fake.Uint32()),
	}
}

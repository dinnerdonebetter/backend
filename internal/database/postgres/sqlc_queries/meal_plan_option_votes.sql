-- name: ArchiveMealPlanOptionVote :execrows

UPDATE meal_plan_option_votes SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_meal_plan_option = $1 AND id = $2;

-- name: CreateMealPlanOptionVote :exec

INSERT INTO meal_plan_option_votes (id,rank,abstain,notes,by_user,belongs_to_meal_plan_option) VALUES ($1,$2,$3,$4,$5,$6);

-- name: CheckMealPlanOptionVoteExistence :one

SELECT EXISTS (
	SELECT
	 meal_plan_option_votes.id
	FROM
	 meal_plan_option_votes
		JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)
	AND meal_plan_option_votes.id = sqlc.arg(meal_plan_option_vote_id)
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = sqlc.arg(meal_plan_event_id)
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
	AND meal_plan_options.id = sqlc.arg(meal_plan_option_id)
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = sqlc.arg(meal_plan_id)
);

-- name: GetMealPlanOptionVotesForMealPlanOption :many

SELECT
    meal_plan_option_votes.id,
    meal_plan_option_votes.rank,
    meal_plan_option_votes.abstain,
    meal_plan_option_votes.notes,
    meal_plan_option_votes.by_user,
    meal_plan_option_votes.created_at,
    meal_plan_option_votes.last_updated_at,
    meal_plan_option_votes.archived_at,
    meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plan_option_votes
         JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
         JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
         JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
  AND meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)
  AND meal_plan_options.archived_at IS NULL
  AND meal_plan_options.belongs_to_meal_plan_event = sqlc.arg(meal_plan_event_id)
  AND meal_plan_options.id = sqlc.arg(meal_plan_option_id)
  AND meal_plan_events.archived_at IS NULL
  AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
  AND meal_plan_events.id = sqlc.arg(meal_plan_event_id)
  AND meal_plans.archived_at IS NULL
  AND meal_plans.id = sqlc.arg(meal_plan_id);

-- name: GetMealPlanOptionVotes :many

SELECT
    meal_plan_option_votes.id,
    meal_plan_option_votes.rank,
    meal_plan_option_votes.abstain,
    meal_plan_option_votes.notes,
    meal_plan_option_votes.by_user,
    meal_plan_option_votes.created_at,
    meal_plan_option_votes.last_updated_at,
    meal_plan_option_votes.archived_at,
    meal_plan_option_votes.belongs_to_meal_plan_option,
    (
        SELECT
            COUNT(meal_plan_events.id)
        FROM
            meal_plan_option_votes
        WHERE
            meal_plan_option_votes.archived_at IS NULL
            AND meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)
            AND meal_plan_option_votes.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND meal_plan_option_votes.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (
                meal_plan_option_votes.last_updated_at IS NULL
                OR meal_plan_option_votes.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
            AND (
                meal_plan_option_votes.last_updated_at IS NULL
                OR meal_plan_option_votes.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
    ) AS filtered_count,
    (
        SELECT
            COUNT(meal_plan_option_votes.id)
        FROM
            meal_plan_option_votes
        WHERE
            meal_plan_option_votes.archived_at IS NULL
    ) AS total_count
FROM meal_plan_option_votes
    JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
    JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
    JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
    AND meal_plan_option_votes.belongs_to_meal_plan_option = sqlc.arg(meal_plan_option_id)
    AND meal_plan_options.archived_at IS NULL
    AND meal_plan_options.belongs_to_meal_plan_event = sqlc.arg(meal_plan_event_id)
    AND meal_plan_options.id = sqlc.arg(meal_plan_option_id)
    AND meal_plan_events.archived_at IS NULL
    AND meal_plan_events.belongs_to_meal_plan = sqlc.arg(meal_plan_id)
    AND meal_plan_events.id = sqlc.arg(meal_plan_event_id)
    AND meal_plans.archived_at IS NULL
    AND meal_plans.id = sqlc.arg(meal_plan_id)
    AND meal_plan_option_votes.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND meal_plan_option_votes.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (
        meal_plan_option_votes.last_updated_at IS NULL
        OR meal_plan_option_votes.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
    AND (
        meal_plan_option_votes.last_updated_at IS NULL
        OR meal_plan_option_votes.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
GROUP BY
    meal_plan_option_votes.id,
    meal_plan_options.id,
    meal_plan_events.id,
    meal_plans.id
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetMealPlanOptionVote :one

SELECT
	meal_plan_option_votes.id,
	meal_plan_option_votes.rank,
	meal_plan_option_votes.abstain,
	meal_plan_option_votes.notes,
	meal_plan_option_votes.by_user,
	meal_plan_option_votes.created_at,
	meal_plan_option_votes.last_updated_at,
	meal_plan_option_votes.archived_at,
	meal_plan_option_votes.belongs_to_meal_plan_option
FROM meal_plan_option_votes
	JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = $1
	AND meal_plan_option_votes.id = $2
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $3
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $4
	AND meal_plan_options.id = $1
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $4;

-- name: UpdateMealPlanOptionVote :execrows

UPDATE meal_plan_option_votes SET rank = $1, abstain = $2, notes = $3, by_user = $4, last_updated_at = NOW() WHERE archived_at IS NULL AND belongs_to_meal_plan_option = $5 AND id = $6;

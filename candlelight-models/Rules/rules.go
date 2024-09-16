package Rules

import "candlelight-models/Actions"

//A set of "Rules" for the game, which is admittedly a loose name. A Rule currently
//might be a constraint or an event-action pair. See their respective files for what each type entails
type RuleSet struct {
	Id               string            `json:"id"`
	Constraints      []Constraint      `json:"constraints"`
	EventActionPairs []EventActionPair `json:"eventActionPairs"`
}

type GameRule struct {
	Id string `json:"id"`
}

//An Event-Action Pair is used to define Actions that should happen on their own
//after a certain Action occurs
type EventActionPair struct {
	GameRule
	Trigger Actions.ActionSet `json:"triggers"`
	Actions Actions.ActionSet `json:"actions"`
}

//A constraint is used to define what a Player can/can't do during certain points
//of the Game.
type Constraint struct {
	GameRule
	//Should be a list Actions currently allowed while the constraint is active
	AllowedActions Actions.ActionSet `json:"allowedActions"`
	//If true, treat AllowedActions as a blacklist of actions instead of a whitelist
	FlipActionsArray bool `json:"flipActionsArray"`
	//A list of Phases this constraint should be active
	ActivePhases []GamePhase `json:"activePhases"`
	//If true, treat the ActivePhases as a list of phases this rule should NOT be active
	FlipActivePhases bool `json:"flipActivePhases"`
}

//A user-defined modifier for when certain actions are allowed
type ActionModifier struct {
	GameRule
	//Actions that are to have the ConstraintFunction applied to
	AffectedActions Actions.ActionSet `json:"affectedActions"`
	//If TRUE, treat the AffectedActions as a set of Actions that should NOT have the ConstraintFunction applied i.e. are exempted
	FlipAffectedActions bool `json:"flipAffectedActions"`
	//A string of valid javascript that can be evaluated to true or false. True means the AffectedActions are allowed, false means they are not
	ConstraintFunction string `json:"constraintFunction"`
}

//Designed to give the system a way of tracking what point of the game we're in. We'll likely
//have some hard-coded ones in place for the creation process, such as Game Start, etc.
type GamePhase struct {
	//Id for book keeping
	Id string `json:"id"`
	//Name of this Phase
	Name string `json:"name"`
}

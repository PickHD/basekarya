package constants

type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "AVAILABLE"
	AssetStatusAssigned    AssetStatus = "ASSIGNED"
	AssetStatusMaintenance AssetStatus = "MAINTENANCE"
	AssetStatusRetired     AssetStatus = "RETIRED"
)

type AssetCondition string

const (
	AssetConditionGood    AssetCondition = "GOOD"
	AssetConditionFair    AssetCondition = "FAIR"
	AssetConditionDamaged AssetCondition = "DAMAGED"
	AssetConditionLost    AssetCondition = "LOST"
)

type AssetAssignmentStatus string

const (
	AssetAssignmentStatusPending  AssetAssignmentStatus = "PENDING"
	AssetAssignmentStatusActive   AssetAssignmentStatus = "ACTIVE"
	AssetAssignmentStatusReturned AssetAssignmentStatus = "RETURNED"
	AssetAssignmentStatusRejected AssetAssignmentStatus = "REJECTED"
)

type AssetAssignmentAction string

const (
	AssetAssignmentActionApprove AssetAssignmentAction = "APPROVE"
	AssetAssignmentActionReject  AssetAssignmentAction = "REJECT"
)

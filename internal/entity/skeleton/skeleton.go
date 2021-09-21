package skeleton

// Skeleton model
type Skeleton struct {
	SkeletonID   int `db:"skeleton_id" json:"skeleton_id"`
	SkeletonName int `db:"skeleton_name" json:"skeleton_name"`
}

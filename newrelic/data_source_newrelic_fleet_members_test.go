//go:build integration || FLEET

package newrelic

// Acceptance tests for the newrelic_fleet_members data source have been
// integrated into TestAccNewRelicFleetMembers_Lifecycle in
// resource_newrelic_fleet_members_test.go.
//
// The standalone tests that previously lived here each created their own fleet
// and added a single entity to it just to read it back — simple enough, but
// they required a separate entity pool to avoid parallel conflicts with the
// resource lifecycle test (which already owns the same entity IDs). Keeping
// them separate would have also meant the data source was only ever verified
// in an almost-empty, single-step scenario.
//
// The integrated approach is better on both counts:
//   - No entity pool conflicts: the data source steps re-use the same fleet
//     and entities that the lifecycle test is already managing, so no
//     additional environment setup is needed.
//   - Richer coverage: the data source is verified at two meaningful points in
//     the lifecycle — after a single-ring create (exercises the unfiltered
//     all-members mode) and after a multi-ring setup (exercises both the
//     unfiltered mode and the ring-filter mode side-by-side), confirming the
//     data source accurately reflects what the resource actually applied.

package metadata

import (
	"errors"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/root-gg/skybook/common"
)

// testBackend creates a temporary database for testing.
func testBackend(t *testing.T) *Backend {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	config := common.DatabaseConfig{Path: dbPath}
	backend, err := NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}
	t.Cleanup(func() { backend.Shutdown() })
	return backend
}

// testJump creates a minimal valid jump for testing.
func testJump(userID uint) *common.Jump {
	return &common.Jump{
		UserID:   userID,
		Date:     common.Today(),
		Dropzone: "Test DZ",
		JumpType: common.JumpTypeFF,
	}
}

func TestAnonymousUserSeeded(t *testing.T) {
	b := testBackend(t)

	var user common.User
	err := b.DB().First(&user, 1).Error
	if err != nil {
		t.Fatalf("anonymous user not found: %v", err)
	}
	if user.Provider != "local" {
		t.Errorf("expected provider 'local', got %q", user.Provider)
	}
	if user.Name != "Skydiver" {
		t.Errorf("expected name 'Skydiver', got %q", user.Name)
	}
}

func TestCreateJump_Append(t *testing.T) {
	b := testBackend(t)

	j1 := testJump(1)
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create jump 1: %v", err)
	}
	if j1.Number != 1 {
		t.Errorf("expected Number 1, got %d", j1.Number)
	}

	j2 := testJump(1)
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("create jump 2: %v", err)
	}
	if j2.Number != 2 {
		t.Errorf("expected Number 2, got %d", j2.Number)
	}

	j3 := testJump(1)
	if err := b.CreateJump(j3); err != nil {
		t.Fatalf("create jump 3: %v", err)
	}
	if j3.Number != 3 {
		t.Errorf("expected Number 3, got %d", j3.Number)
	}
}

func TestInsertJumpAt_Start(t *testing.T) {
	b := testBackend(t)

	// Create 3 jumps
	for i := 0; i < 3; i++ {
		j := testJump(1)
		b.CreateJump(j)
	}

	// Insert at position 1
	newJump := testJump(1)
	newJump.Dropzone = "Inserted"
	if err := b.InsertJumpAt(newJump, 1); err != nil {
		t.Fatalf("insert at 1: %v", err)
	}

	// Verify the inserted jump is #1
	got, err := b.GetJumpByNumber(1, 1)
	if err != nil {
		t.Fatalf("get jump #1: %v", err)
	}
	if got.Dropzone != "Inserted" {
		t.Errorf("expected inserted jump at #1, got dropzone %q", got.Dropzone)
	}

	// Verify total count is now 4
	count, _ := b.CountJumps(1)
	if count != 4 {
		t.Errorf("expected 4 jumps, got %d", count)
	}

	// Verify contiguity
	assertContiguous(t, b, 1, 4)
}

func TestInsertJumpAt_Middle(t *testing.T) {
	b := testBackend(t)

	// Create 5 jumps
	for i := 0; i < 5; i++ {
		j := testJump(1)
		b.CreateJump(j)
	}

	// Insert at position 3
	newJump := testJump(1)
	newJump.Dropzone = "Middle"
	if err := b.InsertJumpAt(newJump, 3); err != nil {
		t.Fatalf("insert at 3: %v", err)
	}

	got, _ := b.GetJumpByNumber(1, 3)
	if got.Dropzone != "Middle" {
		t.Errorf("expected inserted jump at #3")
	}

	count, _ := b.CountJumps(1)
	if count != 6 {
		t.Errorf("expected 6 jumps, got %d", count)
	}

	assertContiguous(t, b, 1, 6)
}

func TestInsertJumpAt_End(t *testing.T) {
	b := testBackend(t)

	// Create 3 jumps
	for i := 0; i < 3; i++ {
		b.CreateJump(testJump(1))
	}

	// Insert at position 4 (append)
	newJump := testJump(1)
	if err := b.InsertJumpAt(newJump, 4); err != nil {
		t.Fatalf("insert at 4: %v", err)
	}

	count, _ := b.CountJumps(1)
	if count != 4 {
		t.Errorf("expected 4 jumps, got %d", count)
	}

	assertContiguous(t, b, 1, 4)
}

func TestInsertJumpAt_OutOfRange(t *testing.T) {
	b := testBackend(t)

	b.CreateJump(testJump(1))

	// Position 3 is out of range for a 1-jump logbook (max is 2)
	err := b.InsertJumpAt(testJump(1), 3)
	if err == nil {
		t.Error("expected error for out of range position")
	}
}

func TestDeleteJump_Last(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 3; i++ {
		b.CreateJump(testJump(1))
	}

	// Delete last jump (#3)
	j, _ := b.GetJumpByNumber(1, 3)
	if err := b.DeleteJump(j); err != nil {
		t.Fatalf("delete: %v", err)
	}

	count, _ := b.CountJumps(1)
	if count != 2 {
		t.Errorf("expected 2 jumps, got %d", count)
	}

	assertContiguous(t, b, 1, 2)
}

func TestDeleteJump_First(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 3; i++ {
		b.CreateJump(testJump(1))
	}

	// Delete first jump (#1)
	j, _ := b.GetJumpByNumber(1, 1)
	if err := b.DeleteJump(j); err != nil {
		t.Fatalf("delete: %v", err)
	}

	count, _ := b.CountJumps(1)
	if count != 2 {
		t.Errorf("expected 2 jumps, got %d", count)
	}

	// Verify remaining jumps are #1 and #2 (renumbered)
	assertContiguous(t, b, 1, 2)
}

func TestDeleteJump_Middle(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 5; i++ {
		b.CreateJump(testJump(1))
	}

	// Delete jump #3
	j, _ := b.GetJumpByNumber(1, 3)
	if err := b.DeleteJump(j); err != nil {
		t.Fatalf("delete: %v", err)
	}

	count, _ := b.CountJumps(1)
	if count != 4 {
		t.Errorf("expected 4 jumps, got %d", count)
	}

	assertContiguous(t, b, 1, 4)
}

func TestEmptyLogbook(t *testing.T) {
	b := testBackend(t)

	count, _ := b.CountJumps(1)
	if count != 0 {
		t.Errorf("expected 0 jumps, got %d", count)
	}

	// Append to empty logbook
	j := testJump(1)
	if err := b.CreateJump(j); err != nil {
		t.Fatalf("create: %v", err)
	}
	if j.Number != 1 {
		t.Errorf("expected Number 1, got %d", j.Number)
	}
}

func TestSingleJump_Delete(t *testing.T) {
	b := testBackend(t)

	j := testJump(1)
	b.CreateJump(j)

	if err := b.DeleteJump(j); err != nil {
		t.Fatalf("delete: %v", err)
	}

	count, _ := b.CountJumps(1)
	if count != 0 {
		t.Errorf("expected 0 jumps after delete, got %d", count)
	}
}

func TestMultiUser_Isolation(t *testing.T) {
	b := testBackend(t)

	// Create a second user
	user2 := &common.User{Provider: "local", ProviderID: "2", Name: "User2"}
	b.DB().Create(user2)

	// Each user gets independent numbering
	b.CreateJump(testJump(1))
	b.CreateJump(testJump(1))
	b.CreateJump(testJump(user2.ID))

	count1, _ := b.CountJumps(1)
	count2, _ := b.CountJumps(user2.ID)

	if count1 != 2 {
		t.Errorf("user 1: expected 2 jumps, got %d", count1)
	}
	if count2 != 1 {
		t.Errorf("user 2: expected 1 jump, got %d", count2)
	}

	// User 2's jump should be #1 (independent sequence)
	j, _ := b.GetJumpByNumber(user2.ID, 1)
	if j == nil {
		t.Fatal("user 2 jump #1 not found")
	}
}

func TestGetJumps_Pagination(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 10; i++ {
		b.CreateJump(testJump(1))
	}

	jumps, total, err := b.GetJumps(1, 0, 5, "number", "asc", JumpFilters{})
	if err != nil {
		t.Fatalf("GetJumps: %v", err)
	}
	if total != 10 {
		t.Errorf("expected total 10, got %d", total)
	}
	if len(jumps) != 5 {
		t.Errorf("expected 5 results, got %d", len(jumps))
	}
	if jumps[0].Number != 1 {
		t.Errorf("expected first jump #1, got #%d", jumps[0].Number)
	}
}

func TestMoveJump_Up(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 3; i++ {
		b.CreateJump(testJump(1))
	}

	j3, _ := b.GetJumpByNumber(1, 3)
	if err := b.MoveJump(j3, 1); err != nil {
		t.Fatalf("MoveJump(3→1): %v", err)
	}
	assertContiguous(t, b, 1, 3)

	// The moved jump should now be at #1
	j, _ := b.GetJump(1, j3.ID)
	if j.Number != 1 {
		t.Errorf("expected moved jump at #1, got %d", j.Number)
	}
}

func TestMoveJump_Down(t *testing.T) {
	b := testBackend(t)

	for i := 0; i < 3; i++ {
		b.CreateJump(testJump(1))
	}

	j1, _ := b.GetJumpByNumber(1, 1)
	if err := b.MoveJump(j1, 3); err != nil {
		t.Fatalf("MoveJump(1→3): %v", err)
	}
	assertContiguous(t, b, 1, 3)

	j, _ := b.GetJump(1, j1.ID)
	if j.Number != 3 {
		t.Errorf("expected moved jump at #3, got %d", j.Number)
	}
}

func assertContiguous(t *testing.T, b *Backend, userID uint, expectedCount int64) {
	t.Helper()
	for i := uint(1); i <= uint(expectedCount); i++ {
		j, err := b.GetJumpByNumber(userID, i)
		if err != nil {
			t.Errorf("jump #%d not found: %v", i, err)
			continue
		}
		if j.Number != i {
			t.Errorf("expected jump #%d, got #%d", i, j.Number)
		}
	}
}

// TestUpdateJump_CreatedAtUnchanged verifies that UpdateJump never overwrites created_at.
func TestUpdateJump_CreatedAtUnchanged(t *testing.T) {
	b := testBackend(t)

	j := testJump(1)
	if err := b.CreateJump(j); err != nil {
		t.Fatalf("create: %v", err)
	}
	createdAt := j.CreatedAt

	// Simulate a client sending a different created_at — should be ignored.
	j.Dropzone = "Updated DZ"
	if err := b.UpdateJump(j); err != nil {
		t.Fatalf("update: %v", err)
	}

	// Re-fetch from DB to confirm.
	fetched, err := b.GetJump(1, j.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !fetched.CreatedAt.Equal(createdAt) {
		t.Errorf("created_at changed: was %v, now %v", createdAt, fetched.CreatedAt)
	}
	if fetched.Dropzone != "Updated DZ" {
		t.Errorf("dropzone not updated: got %q", fetched.Dropzone)
	}
}

// TestUpdateJump_BooleanOverwrite verifies PUT semantics: a false boolean in the
// body overwrites an existing true — the frontend is expected to send all fields.
func TestUpdateJump_BooleanOverwrite(t *testing.T) {
	b := testBackend(t)

	// Create with NightJump=true.
	j := &common.Jump{
		UserID:    1,
		Date:      common.Today(),
		Dropzone:  "DZ",
		JumpType:  common.JumpTypeFF,
		NightJump: true,
	}
	if err := b.CreateJump(j); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Update with NightJump=false (full-body PUT — frontend explicitly cleared it).
	j.NightJump = false
	j.Dropzone = "DZ Updated"
	if err := b.UpdateJump(j); err != nil {
		t.Fatalf("update: %v", err)
	}

	fetched, err := b.GetJump(1, j.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if fetched.NightJump {
		t.Error("NightJump should be false after PUT with nightJump=false")
	}
	if fetched.Dropzone != "DZ Updated" {
		t.Errorf("dropzone not updated: got %q", fetched.Dropzone)
	}
}

// dateOf creates a DateOnly for a given year, month, day (UTC midnight).
func dateOf(year, month, day int) common.DateOnly {
	return common.NewDateOnly(year, time.Month(month), day)
}

// TestFirstJump_NoValidation — first jump in empty logbook always passes.
func TestFirstJump_NoValidation(t *testing.T) {
	b := testBackend(t)
	j := testJump(1)
	if err := b.CreateJump(j); err != nil {
		t.Fatalf("first jump should always succeed: %v", err)
	}
}

// TestCreateJump_SameDay_Valid — two jumps on the same date is allowed.
func TestCreateJump_SameDay_Valid(t *testing.T) {
	b := testBackend(t)
	today := dateOf(2025, 6, 15)

	j1 := &common.Jump{UserID: 1, Date: today, Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	j2 := &common.Jump{UserID: 1, Date: today, Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("two jumps on the same day should succeed: %v", err)
	}
}

// TestCreateJump_DateOrder_AppendValid — appending with date ≥ last passes.
func TestCreateJump_DateOrder_AppendValid(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 15), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("append with later date should succeed: %v", err)
	}
}

// TestCreateJump_DateOrder_AppendInvalid — appending with date < last returns DateOrderError.
func TestCreateJump_DateOrder_AppendInvalid(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 15), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	err := b.CreateJump(j2)
	if err == nil {
		t.Fatal("expected DateOrderError, got nil")
	}
	var doe *common.DateOrderError
	if !errors.As(err, &doe) {
		t.Errorf("expected *DateOrderError, got %T: %v", err, err)
	}
}

// TestInsertJumpAt_DateOrder_Valid — inserting between two jumps with a fitting date passes.
func TestInsertJumpAt_DateOrder_Valid(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 15), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("create j2: %v", err)
	}

	// Insert at position 2 (between j1 and j2) with a date squarely in between
	ins := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 8), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.InsertJumpAt(ins, 2); err != nil {
		t.Fatalf("insert with date between neighbors should succeed: %v", err)
	}
}

// TestInsertJumpAt_DateOrder_InvalidBefore — inserting with date < prev returns DateOrderError.
func TestInsertJumpAt_DateOrder_InvalidBefore(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 20), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("create j2: %v", err)
	}

	// Insert at position 2 but with date before j1 (June1 < June10)
	ins := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	err := b.InsertJumpAt(ins, 2)
	if err == nil {
		t.Fatal("expected DateOrderError (date before prev), got nil")
	}
	var doe *common.DateOrderError
	if !errors.As(err, &doe) {
		t.Errorf("expected *DateOrderError, got %T: %v", err, err)
	}
}

// TestInsertJumpAt_DateOrder_InvalidAfter — inserting with date > next returns DateOrderError.
func TestInsertJumpAt_DateOrder_InvalidAfter(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 20), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	if err := b.CreateJump(j1); err != nil {
		t.Fatalf("create j1: %v", err)
	}
	if err := b.CreateJump(j2); err != nil {
		t.Fatalf("create j2: %v", err)
	}

	// Insert at position 2 but with date after j2 (June25 > June20)
	ins := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 25), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	err := b.InsertJumpAt(ins, 2)
	if err == nil {
		t.Fatal("expected DateOrderError (date after next), got nil")
	}
	var doe *common.DateOrderError
	if !errors.As(err, &doe) {
		t.Errorf("expected *DateOrderError, got %T: %v", err, err)
	}
}

// TestUpdateJump_DateOrder_Valid — changing date within bounds succeeds.
func TestUpdateJump_DateOrder_Valid(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 15), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j3 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 30), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	for _, j := range []*common.Jump{j1, j2, j3} {
		if err := b.CreateJump(j); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	// Move j2's date; still between j1 and j3
	j2.Date = dateOf(2025, 6, 20)
	if err := b.UpdateJump(j2); err != nil {
		t.Fatalf("update within bounds should succeed: %v", err)
	}
}

// TestUpdateJump_DateOrder_Invalid — changing date outside bounds returns DateOrderError.
func TestUpdateJump_DateOrder_Invalid(t *testing.T) {
	b := testBackend(t)

	j1 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j2 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 15), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	j3 := &common.Jump{UserID: 1, Date: dateOf(2025, 6, 30), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	for _, j := range []*common.Jump{j1, j2, j3} {
		if err := b.CreateJump(j); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	// Try to push j2 after j3 — should fail
	j2.Date = dateOf(2025, 7, 10)
	err := b.UpdateJump(j2)
	if err == nil {
		t.Fatal("expected DateOrderError (date after next), got nil")
	}
	var doe *common.DateOrderError
	if !errors.As(err, &doe) {
		t.Errorf("expected *DateOrderError, got %T: %v", err, err)
	}
}

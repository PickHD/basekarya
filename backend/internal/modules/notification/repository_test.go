package notification

import (
	"testing"
	"time"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupNotificationTestDB(t *testing.T) *testutil.TestDB {
	t.Helper()
	tdb := testutil.NewTestDB(
		&user.User{},
		&Notification{},
	)
	t.Cleanup(tdb.Close)
	return tdb
}

func seedNotificationTestData(t *testing.T, db *testutil.TestDB) {
	t.Helper()
	companyID := uint(1)

	require.NoError(t, db.DB.Create(&user.User{ID: 1, Username: "john", CompanyID: companyID}).Error)
	require.NoError(t, db.DB.Create(&Notification{
		CompanyID: companyID, UserID: 1, Type: "LEAVE", Title: "New Leave",
		Message: "You have a new leave request", RelatedID: 1, IsRead: false,
	}).Error)
	require.NoError(t, db.DB.Create(&Notification{
		CompanyID: companyID, UserID: 1, Type: "LEAVE", Title: "Approved",
		Message: "Your leave was approved", RelatedID: 2, IsRead: true, CreatedAt: time.Now().AddDate(0, 0, -10),
	}).Error)
}

func TestRepoNotif_Create(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	notif := &Notification{
		CompanyID: 1, UserID: 1, Type: "OVERTIME", Title: "New Overtime",
		Message: "You have overtime request", RelatedID: 5, IsRead: false,
	}

	err := repo.Create(ctx, notif)
	require.NoError(t, err)
	assert.NotZero(t, notif.ID)
}

func TestRepoNotif_FindByID(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{name: "success", id: 1, wantErr: false},
		{name: "not found", id: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif, err := repo.FindByID(ctx, tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, notif.ID)
				assert.Equal(t, "New Leave", notif.Title)
			}
		})
	}
}

func TestRepoNotif_FindAllByUserID(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	notifs, err := repo.FindAllByUserID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, notifs, 2)
}

func TestRepoNotif_MarkAsRead(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	err := repo.MarkAsRead(ctx, 1)
	require.NoError(t, err)

	notif, _ := repo.FindByID(ctx, 1)
	assert.True(t, notif.IsRead)
}

func TestRepoNotif_DeleteReadOlderThan(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	err := repo.DeleteReadOlderThan(ctx, 5)
	require.NoError(t, err)

	notifs, _ := repo.FindAllByUserID(ctx, 1)
	assert.Len(t, notifs, 1)
}

func TestRepoNotif_CreateBatch(t *testing.T) {
	tdb := setupNotificationTestDB(t)
	repo := NewRepository(tdb.DB)
	ctx := testutil.CtxWithTenant(1, 1, false)

	seedNotificationTestData(t, tdb)

	batch := []*Notification{
		{CompanyID: 1, UserID: 1, Type: "LEAVE", Title: "N1", Message: "M1", IsRead: false},
		{CompanyID: 1, UserID: 1, Type: "LEAVE", Title: "N2", Message: "M2", IsRead: false},
	}

	err := repo.CreateBatch(ctx, batch)
	require.NoError(t, err)

	notifs, _ := repo.FindAllByUserID(ctx, 1)
	assert.Len(t, notifs, 4)
}

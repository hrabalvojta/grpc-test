package server

import (
	"context"
	"io/ioutil"
	"os"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx/v4/pgxpool"

	pbExample "github.com/hrabalvojta/grpc-test/proto"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
	users []*pbExample.User
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, req *pbExample.AddUserRequest) (*pbExample.User, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	user := &pbExample.User{
		//Id:    uuid.Must(uuid.NewV4()).String(),
		//Id:    uuid.Must(uuid.NewV4()).String(),
		Email: req.GetEmail(),
	}
	b.users = append(b.users, user)

	return user, nil
}

// GetUser from the in-memory store.
func (b *Backend) GetUser(ctx context.Context, req *pbExample.GetUserRequest) (*pbExample.User, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	//SELECT customer_id, first_name, last_name, email FROM public.customer WHERE customer_id=100;

	for _, user := range b.users {
		if user.Id == req.GetId() {
			return user, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "user with ID %q could not be found", req.GetId())
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *pbExample.ListUsersRequest, srv pbExample.UserService_ListUsersServer) error {

	// init variables functions
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
	ctx := context.Background()
	pgDbUrl := "postgres://postgres:secret@127.0.0.1:5432/dvdrental"
	pgQuery := `SELECT customer_id, first_name, last_name, email FROM public.customer ORDER BY 1;`

	// Create connection pool
	dbpool, err := pgxpool.Connect(ctx, pgDbUrl)
	if err != nil {
		return status.Errorf(codes.Unavailable, "Unable to connect to database %s", err)
	}
	defer dbpool.Close()
	log.Info("[sql-query] Finished connection")

	// Run query
	rows, err := dbpool.Query(ctx, pgQuery)
	if err != nil {
		return status.Errorf(codes.Unavailable, "Unable to generate data: %s", err)
	}
	defer rows.Close()
	log.Info("[sql-query] Finished get data")

	// Scan row line by line and Send
	for rows.Next() {
		var user pbExample.User
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return status.Errorf(codes.Unavailable, "Unable to scan %v\n", err)
		}
		log.Info("[sql-query] ID: ", user.Id, "; First name: ", user.FirstName, "; Lastname: ", user.LastName, "; Email: ", user.Email)
		err = srv.Send(&user)
		if err != nil {
			return err
		}
	}

	return nil
}

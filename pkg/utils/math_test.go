package utils

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/* Common Testing functions & Methods from "testing" module
- t.Errorf(format string, args ...) & t.Error(args ...) : logs an error message but continue with the test
- t.Fatalf(format string, args ...) & t.Fatal(args ...) : logs an error message and stops the execution of current test
- t.Skipf(format string, args ...)  & t.Skip(args ...)  : Useful for skipping the test baed on conditions
- t.Run(name string, f func(t *testing.T))              : Allows you to create subtests. Useful for table-driven testing.

- t.Parallel()                 : Marks the test as being able to run in parallel with other parallel tests.
- t.Helper()                   : Marks the calling function as a test helper function. When printing file and line information, that function is skipped.
                                 This is useful for custom assertion functions.
- t.TempDir()                  : Creates temporary directory
- TestMain(m *testing.M)       : the testing package provides a TestMain function that allows you to implement setup and teardown logic
- Benchmark<func>(b*testing.B) : Built in performance benchmark, must start with Benchmark (b.N is automatically adjusted)
- Test Suites with t.Run       : Go does not have native support for test suites, but you can use t.Run to organize code into Subtests
- testify/mock & testify/assert: create mocks, stubs, fakes and spies
  testify/assert and require: write more easy & readable tests.
	  assert.Equal(t, expected, actual, "Optional format string")
	  assert.Nil(t, object, "Expected object to be nil")
	  assert.NotNil(t, result, "result should not be nil")
	  assert.True(t, condition, "Expected condition to be true")
  	  assert.False(t, condition, "Expected condition to be false")
	require assertion fails, the test will immediately stop executing
	  require.Equal(t, expected, actual, "Optional format string")

  Mocks: Ensure methods are called with expected arguments and can return predefined responses;
  	     primarily used for verifying interactions.
  Spies: Record method calls and arguments for later inspection;
         verification is performed manually by inspecting recorded data.
  Stubs: Provide simple, predefined responses to simulate controlled behavior;
         used for basic isolation of code under test.
  Fakes: Offer more complex, working implementations, often with in-memory state;
         used to simulate real components closely but not suitable for production.

go test ./...           # Test all packages
go test -v              # Run tests with verbose output
go test -run TestSum    # Run a specific test
go test -bench=.        # Run all benchmarks
go test -coverprofile=coverage.out && go tool cover -html=coverage.out   # Generate an HTML coverage report

*/

func helperFunction(t *testing.T) {
	t.Helper() // marks this function as a helper
	// ...
}

/** test setup & tear down **/

func globalSetup() {
	// Global setup actions
}

func globalTeardown() {
	// Global teardown actions
}

func setupTest(t *testing.T) func() {
	// Setup for individual tests
	t.Helper()         // Mark this function as a helper
	dir := t.TempDir() // Example: create a temporary directory
	fmt.Println(dir)
	// return a cleanup function to be deferred
	return func() {
		// Teardown actions: remove temporary files, etc.
		// Example: os.RemoveAll(dir)
	}
}

func TestMain(m *testing.M) {
	globalSetup()
	code := m.Run()
	globalTeardown()
	os.Exit(code)
}

func TestSomeFunction(t *testing.T) {
	teardown := setupTest(t) // Call setup
	defer teardown()         // Ensure teardown is called after test

	// Your test logic goes here
}

/** Per test setup & tear down **/

func TestSubstract(t *testing.T) {
	t.Parallel()
	expected := 1

	if result := Substract(3, 2); result != expected {
		t.Errorf("Substract(2, 3) = %d; expected %d", result, expected)
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		a, b, sum int
	}{
		{"1+1", 1, 1, 2},
		{"2+2", 2, 2, 4},
		{"2+3", 2, 3, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.sum {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.sum)
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	// The run loop for the benchmark
	for i := 0; i < b.N; i++ {
		_ = Add(3, 5)
	}
}

func BenchmarkSumWithSetup(b *testing.B) {
	// Setup (not timed)
	a, c := 3, 5

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_ = Add(a, c)
	}
}

// Test Suites - with nested tests using t.Run
func TestArithmetic(t *testing.T) {
	t.Run("Sum", func(t *testing.T) {
		if result := Add(3, 5); result != 8 {
			t.Errorf("Sum(3, 5) = %d; want 8", result)
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		if result := Substract(5, 3); result != 2 {
			t.Errorf("Subtract(5, 3) = %d; want 2", result)
		}
	})
}

// Assert & require with testify
func TestMath(t *testing.T) {
	result := Add(2, 3)
	assert.Equal(t, 5, result, "they should be equal")

	result1 := Add(2, 3)
	require.NotNil(t, result1, "Expected object not to be nil")
	// Further assertions will not be executed if any require fails
}

/************* Explain Mock ******************/
/* Mocks, Spies, Stubs & Fakes
myapp/
├── user.go             // Contains the User struct and UserRepository interface
├── user_service.go     // Contains the UserService with operations on users
├── mock_user_repo.go   // Contains the MockUserRepository implementation
└── user_service_test.go// Contains the tests for UserService */

// user.go
type User struct{ ID, Name string }

type UserRepository interface {
	GetUserByID(id string) (*User, error)
	SaveUser(user *User) error
}

// user_service.go
type UserService struct{ Repo UserRepository }

func (s *UserService) GetUser(id string) (*User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) CreateUser(id string, name string) error {
	user := &User{ID: id, Name: name}
	return s.Repo.SaveUser(user)
}

// mock_user_repo.go
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) GetUserByID(id string) (*User, error) {
	args := m.Called(id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) SaveUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

// user_service_test.go
func TestUserService_GetUser_Mock(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := &UserService{Repo: mockRepo}

	// Setup expected mock behavior
	user := &User{ID: "123", Name: "John Doe"}
	mockRepo.On("GetUserByID", "123").Return(user, nil)

	// Call the method under test
	result, err := userService.GetUser("123")

	// Use assertions to verify the result
	require.NoError(t, err)
	assert.Equal(t, "John Doe", result.Name)

	// Verify that the GetUserByID method was called with the expected argument (Spy behavior)
	mockRepo.AssertCalled(t, "GetUserByID", "123")
}

func TestUserService_CreateUser_Stub(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := &UserService{Repo: mockRepo}

	// Stub the SaveUser method
	mockRepo.On("SaveUser", mock.Anything).Return(nil)

	// Call the method under test
	err := userService.CreateUser("123", "John Doe")

	// Use assertions to verify the result
	require.NoError(t, err)

	// Verify that the SaveUser method was called with expected arguments (Spy behavior)
	mockRepo.AssertCalled(t, "SaveUser", &User{ID: "123", Name: "John Doe"})
}

func TestUserService_CreateUser_Fake(t *testing.T) {
	fakeRepo := &FakeUserRepository{}
	userService := &UserService{Repo: fakeRepo}

	// Call the method under test
	err := userService.CreateUser("123", "John Doe")
	require.NoError(t, err)

	// Check if the user was actually saved
	savedUser, err := userService.GetUser("123")
	require.NoError(t, err)
	assert.Equal(t, "John Doe", savedUser.Name)
}

// FakeUserRepository is a simple, working implementation of UserRepository for testing.
type FakeUserRepository struct {
	users map[string]*User
}

func (f *FakeUserRepository) GetUserByID(id string) (*User, error) {
	if user, ok := f.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (f *FakeUserRepository) SaveUser(user *User) error {
	if f.users == nil {
		f.users = make(map[string]*User)
	}
	f.users[user.ID] = user
	return nil
}

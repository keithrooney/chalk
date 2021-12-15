package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

type Person struct {
	Model
	Firstname string
	Surname   string
}

func TestDatabase(t *testing.T) {

	dsn := fmt.Sprintf("/tmp/%s.test.delivr.sqlite", uuid.NewString())
	factory := SqliteDialectorFactory{}

	database := NewDatabase(dsn, factory)

	if err := database.AutoMigrate(&Person{}); err != nil {
		t.Fatal(err)
	}

	t.Run("TestCreate", func(t *testing.T) {

		person := &Person{
			Firstname: "Michael",
			Surname:   "Cunningham",
		}

		if err := database.Create(person); err != nil {
			t.Fatal(err)
		}

		other := &Person{}
		if err := database.Get(person.ID, other); err != nil {
			t.Fatal(err)
		}

		if person.ID != other.ID {
			t.Fatalf(`%v != %v`, person, other)
		}

	})

	t.Run("TestQuery", func(t *testing.T) {

		people := []Person{
			{
				Firstname: "John",
				Surname:   "Doe",
			},
			{
				Firstname: "John",
				Surname:   "Nemo",
			},
		}

		if err := database.Create(people); err != nil {
			t.Fatal(err)
		}

		others := []Person{}
		if err := database.Query(map[string]interface{}{"firstname": "John"}, &others); err != nil {
			t.Fatal(err)
		}

		if len(others) != 2 {
			t.Fatalf("Expected only two results.")
		}

		more := []Person{}
		person := people[0]
		if err := database.Query(map[string]interface{}{"id": person.ID}, &more); err != nil {
			t.Fatal(err)
		}

		if len(more) != 1 {
			t.Fatalf("Expected only one result.")
		}

	})

	t.Run("TestUpdate", func(t *testing.T) {

		person := &Person{
			Firstname: "Michael",
			Surname:   "Mitchell",
		}

		if err := database.Create(person); err != nil {
			t.Fatal(err)
		}

		person.Firstname = "John"
		person.Surname = "Bloomberg"

		if err := database.Update(person); err != nil {
			t.Fatal(err)
		}

		other := &Person{}
		if err := database.Get(person.ID, other); err != nil {
			t.Fatal(err)
		}

		if other.Firstname != "John" {
			t.Fatalf("Expected field to be updated.")
		}
		if other.Surname != "Bloomberg" {
			t.Fatalf("Expected field to be updated.")
		}

	})

	t.Run("TestDelete", func(t *testing.T) {

		person := &Person{
			Firstname: "Samuel L.",
			Surname:   "Jackson",
		}

		if err := database.Create(person); err != nil {
			t.Fatal(err)
		}

		if err := database.Delete(&Person{Model: Model{ID: person.ID}}); err != nil {
			t.Fatal(err)
		}

		people := []Person{}
		if err := database.Query(map[string]interface{}{"id": person.ID}, &people); err != nil {
			t.Fatal(err)
		}

		if len(people) != 0 {
			t.Fatalf("Expected zero results.")
		}

	})

	t.Run("TestDeleteInvalidId", func(t *testing.T) {
		if err := database.Delete(&Person{Model: Model{ID: 1230123}}); err == nil {
			t.Fatal(err)
		}
	})

	t.Cleanup(func() {
		if err := os.Remove(dsn); err != nil {
			t.Fatal(err)
		}
	})

}

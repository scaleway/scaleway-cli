package anonuuid

import (
	"fmt"
	"testing"
)

const exampleInput string = `VOLUMES_0_SERVER_ID=15573749-c89d-41dd-a655-16e79bed52e0
VOLUMES_0_SERVER_NAME=hello
VOLUMES_0_ID=c245c3cb-3336-4567-ada1-70cb1fe4eefe
VOLUMES_0_SIZE=50000000000
ORGANIZATION=fe1e54e8-d69d-4f7c-a9f1-42069e03da31
TEST=15573749-c89d-41dd-a655-16e79bed52e0`

func ExampleAnonUUID_Sanitize() {
	anonuuid := New()
	fmt.Println(anonuuid.Sanitize(exampleInput))
	// Output:
	// VOLUMES_0_SERVER_ID=00000000-0000-0000-0000-000000000000
	// VOLUMES_0_SERVER_NAME=hello
	// VOLUMES_0_ID=11111111-1111-1111-1111-111111111111
	// VOLUMES_0_SIZE=50000000000
	// ORGANIZATION=22222222-2222-2222-2222-222222222222
	// TEST=00000000-0000-0000-0000-000000000000
}

func TestAnonUUID_Cache(t *testing.T) {
	anonuuid := New()
	if len(anonuuid.cache) != 0 {
		t.Fatalf("anonuuid.cache should be empty")
	}

	anonuuid.Sanitize("hello")
	if len(anonuuid.cache) != 0 {
		t.Fatalf("anonuuid.cache should be empty")
	}

	anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
	if len(anonuuid.cache) != 1 {
		t.Fatalf("anonuuid.cache should contain 1 entry")
	}

	anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
	if len(anonuuid.cache) != 1 {
		t.Fatalf("anonuuid.cache should contain 1 entry")
	}

	anonuuid.Sanitize("hello c245c3cb-3336-4567-ada1-70cb1fe4eefe")
	if len(anonuuid.cache) != 2 {
		t.Fatalf("anonuuid.cache should contain 2 entries")
	}

	anonuuid.Sanitize("hello c245c3cb-3336-4567-ada1-70cb1fe4eefe")
	if len(anonuuid.cache) != 2 {
		t.Fatalf("anonuuid.cache should contain 2 entries")
	}

	anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
	if len(anonuuid.cache) != 2 {
		t.Fatalf("anonuuid.cache should contain 2 entries")
	}
}

func ExampleAnonUUID_FakeUUID() {
	anonuuid := New()
	fmt.Println(anonuuid.FakeUUID("15573749-c89d-41dd-a655-16e79bed52e0"))
	fmt.Println(anonuuid.FakeUUID("15573749-c89d-41dd-a655-16e79bed52e0"))
	fmt.Println(anonuuid.FakeUUID("c245c3cb-3336-4567-ada1-70cb1fe4eefe"))
	fmt.Println(anonuuid.FakeUUID("c245c3cb-3336-4567-ada1-70cb1fe4eefe"))
	fmt.Println(anonuuid.FakeUUID("15573749-c89d-41dd-a655-16e79bed52e0"))
	fmt.Println(anonuuid.FakeUUID("fe1e54e8-d69d-4f7c-a9f1-42069e03da31"))
	// Output:
	// 00000000-0000-0000-0000-000000000000
	// 00000000-0000-0000-0000-000000000000
	// 11111111-1111-1111-1111-111111111111
	// 11111111-1111-1111-1111-111111111111
	// 00000000-0000-0000-0000-000000000000
	// 22222222-2222-2222-2222-222222222222
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkAnonUUID_FakeUUID(b *testing.B) {
	anonuuid := New()
	for i := 0; i < b.N; i++ {
		anonuuid.FakeUUID("15573749-c89d-41dd-a655-16e79bed52e0")
	}
}

func BenchmarkAnonUUID_Sanitize(b *testing.B) {
	anonuuid := New()
	for i := 0; i < b.N; i++ {
		anonuuid.Sanitize("A: 15573749-c89d-41dd-a655-16e79bed52e0, B: c245c3cb-3336-4567-ada1-70cb1fe4eefe, A: 15573749-c89d-41dd-a655-16e79bed52e0")
	}
}

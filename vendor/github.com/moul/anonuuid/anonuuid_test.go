package anonuuid

import (
	"fmt"
	"testing"

	. "github.com/moul/anonuuid/vendor/github.com/smartystreets/goconvey/convey"
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
	// VOLUMES_0_SERVER_ID=00000000-0000-1000-0000-000000000000
	// VOLUMES_0_SERVER_NAME=hello
	// VOLUMES_0_ID=11111111-1111-1111-1111-111111111111
	// VOLUMES_0_SIZE=50000000000
	// ORGANIZATION=22222222-2222-1222-2222-222222222222
	// TEST=00000000-0000-1000-0000-000000000000
}

func TestAnonUUID_cache(t *testing.T) {
	Convey("Testing AnonUUID.cache", t, func() {
		anonuuid := New()
		So(len(anonuuid.cache), ShouldEqual, 0)

		anonuuid.Sanitize("hello")
		So(len(anonuuid.cache), ShouldEqual, 0)

		anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
		So(len(anonuuid.cache), ShouldEqual, 1)

		anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
		So(len(anonuuid.cache), ShouldEqual, 1)

		anonuuid.Sanitize("hello c245c3cb-3336-4567-ada1-70cb1fe4eefe")
		So(len(anonuuid.cache), ShouldEqual, 2)

		anonuuid.Sanitize("hello c245c3cb-3336-4567-ada1-70cb1fe4eefe")
		So(len(anonuuid.cache), ShouldEqual, 2)

		anonuuid.Sanitize("hello 15573749-c89d-41dd-a655-16e79bed52e0")
		So(len(anonuuid.cache), ShouldEqual, 2)
	})
}

func TestAnonUUID_Sanitize(t *testing.T) {
	Convey("Testing AnonUUID.Sanitize", t, func() {
		realuuid1 := "15573749-c89d-41dd-a655-16e79bed52e0"
		realuuid2 := "c245c3cb-3336-4567-ada1-70cb1fe4eefe"

		input1 := "hello"
		input2 := "hello " + realuuid1
		input3 := "hello " + realuuid2
		input4 := fmt.Sprintf("hello %s %s %s", realuuid1, realuuid2, realuuid1)

		Convey("no options", func() {
			anonuuid := New()

			out1 := anonuuid.Sanitize(input1)
			out2 := anonuuid.Sanitize(input2)
			out3 := anonuuid.Sanitize(input2)
			out4 := anonuuid.Sanitize(input3)
			out5 := anonuuid.Sanitize(input3)
			out6 := anonuuid.Sanitize(input2)
			out7 := anonuuid.Sanitize(input1)
			out8 := anonuuid.Sanitize(input4)

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid2)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(out1, ShouldEqual, input1)
			So(out2, ShouldEqual, "hello "+fakeuuid1)
			So(out2, ShouldNotEqual, input2)
			So(out3, ShouldEqual, "hello "+fakeuuid1)
			So(out3, ShouldNotEqual, input2)
			So(out4, ShouldEqual, "hello "+fakeuuid2)
			So(out4, ShouldNotEqual, input3)
			So(out5, ShouldEqual, "hello "+fakeuuid2)
			So(out5, ShouldNotEqual, input3)
			So(out6, ShouldEqual, "hello "+fakeuuid1)
			So(out6, ShouldNotEqual, input2)
			So(out7, ShouldEqual, input1)
			So(out8, ShouldEqual, fmt.Sprintf("hello %s %s %s", fakeuuid1, fakeuuid2, fakeuuid1))
			So(out8, ShouldNotEqual, input4)
		})
		Convey("--prefix=XXX", func() {
			anonuuid := New()

			anonuuid.Prefix = "world"

			out1 := anonuuid.Sanitize(input1)
			out2 := anonuuid.Sanitize(input2)
			out3 := anonuuid.Sanitize(input2)
			out4 := anonuuid.Sanitize(input3)
			out5 := anonuuid.Sanitize(input3)
			out6 := anonuuid.Sanitize(input2)
			out7 := anonuuid.Sanitize(input1)
			out8 := anonuuid.Sanitize(input4)

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid2)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(out1, ShouldEqual, input1)
			So(out2, ShouldContainSubstring, "hello world")
			So(out2, ShouldEqual, "hello "+fakeuuid1)
			So(out2, ShouldNotEqual, input2)
			So(out3, ShouldContainSubstring, "hello world")
			So(out3, ShouldEqual, "hello "+fakeuuid1)
			So(out3, ShouldNotEqual, input2)
			So(out4, ShouldContainSubstring, "hello world")
			So(out4, ShouldEqual, "hello "+fakeuuid2)
			So(out4, ShouldNotEqual, input3)
			So(out5, ShouldContainSubstring, "hello world")
			So(out5, ShouldEqual, "hello "+fakeuuid2)
			So(out5, ShouldNotEqual, input3)
			So(out6, ShouldContainSubstring, "hello world")
			So(out6, ShouldEqual, "hello "+fakeuuid1)
			So(out6, ShouldNotEqual, input2)
			So(out7, ShouldEqual, input1)
			So(out8, ShouldContainSubstring, "hello world")
			So(out8, ShouldEqual, fmt.Sprintf("hello %s %s %s", fakeuuid1, fakeuuid2, fakeuuid1))
			So(out8, ShouldNotEqual, input4)
		})
	})
}

func TestAnonUUID_FakeUUID(t *testing.T) {
	Convey("Testing AnonUUID.FakeUUID", t, func() {
		realuuid1 := "15573749-c89d-41dd-a655-16e79bed52e0"
		realuuid2 := "c245c3cb-3336-4567-ada1-70cb1fe4eefe"

		Convey("no options", func() {
			anonuuid := New()

			expected1 := "00000000-0000-1000-0000-000000000000"
			expected2 := "11111111-1111-1111-1111-111111111111"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=hello", func() {
			anonuuid := New()

			anonuuid.Prefix = "hello"

			expected1 := "hello000-0000-1000-0100-000000000000"
			// FIXME: why this                ^ ?
			expected2 := "hello111-1111-1111-1111-111111111111"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=helloworld", func() {
			anonuuid := New()

			anonuuid.Prefix = "helloworld"

			expected1 := "hellowor-ld00-1000-0000-001000000000"
			// FIXME: why this                      ^ ?
			expected2 := "hellowor-ld11-1111-1111-111111111111"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=@@@", func() {
			anonuuid := New()

			anonuuid.Prefix = "@@@"

			expected1 := "invalidp-refi-1000-0000-000001000000"
			// FIXME: why this                         ^ ?
			expected2 := "invalidp-refi-1111-1111-111111111111"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--suffix=hello", func() {
			anonuuid := New()

			anonuuid.Suffix = "hello"

			expected1 := "00000000-0000-1000-0000-0000000hello"
			expected2 := "11111111-1111-1111-1111-1111111hello"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--suffix=helloworldhello", func() {
			anonuuid := New()

			anonuuid.Suffix = "helloworldhello"

			expected1 := "00000000-0000-1000-0hel-loworldhello"
			expected2 := "11111111-1111-1111-1hel-loworldhello"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--suffix=@@@", func() {
			anonuuid := New()

			anonuuid.Suffix = "@@@"

			expected1 := "00000000-0000-1000-0000-0invalsuffix"
			expected2 := "11111111-1111-1111-1111-1invalsuffix"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=hello --suffix=hello", func() {
			anonuuid := New()

			anonuuid.Prefix = "hello"
			anonuuid.Suffix = "hello"

			expected1 := "hello000-0000-1000-0100-0000000hello"
			// FIXME: why this                ^ ?
			expected2 := "hello111-1111-1111-1111-1111111hello"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=helloworldhello --suffix=helloworldhello", func() {
			anonuuid := New()

			anonuuid.Prefix = "helloworldhello"
			anonuuid.Suffix = "helloworldhello"

			expected1 := "hellowor-ldhe-1lo0-0hel-loworldhello"
			// FIXME: should not do this :)
			expected2 := "hellowor-ldhe-1lo1-1hel-loworldhello"
			// FIXME: should not do this :)

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--prefix=helloworldhello --suffix=@@@", func() {
			anonuuid := New()

			anonuuid.Prefix = "@@@"
			anonuuid.Suffix = "@@@"

			expected1 := "invalidp-refi-1000-0000-0invalsuffix"
			expected2 := "invalidp-refi-1111-1111-1invalsuffix"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--keep-beginning", func() {
			anonuuid := New()

			anonuuid.KeepBeginning = true

			expected1 := "15573749-0000-1000-0000-100000000000"
			expected2 := "c245c3cb-1111-1111-1111-111111111111"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--keep-end", func() {
			anonuuid := New()

			anonuuid.KeepEnd = true

			expected1 := "00000000-0000-1000-0000-16e79bed52e0"
			expected2 := "11111111-1111-1111-1111-70cb1fe4eefe"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--keep-beginning --keep-end", func() {
			anonuuid := New()

			anonuuid.KeepBeginning = true
			anonuuid.KeepEnd = true

			expected1 := "15573749-0000-1000-0000-16e79bed52e0"
			expected2 := "c245c3cb-1111-1111-1111-70cb1fe4eefe"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--hexspeak", func() {
			anonuuid := New()

			anonuuid.Hexspeak = true

			expected1 := "0ff1ce0f-f1ce-1ff1-ce0f-f1ce0ff1ce0f"
			expected2 := "31337313-3731-1373-1337-313373133731"

			fakeuuid1 := anonuuid.FakeUUID(realuuid1)
			fakeuuid2 := anonuuid.FakeUUID(realuuid1)
			fakeuuid3 := anonuuid.FakeUUID(realuuid2)
			fakeuuid4 := anonuuid.FakeUUID(realuuid2)
			fakeuuid5 := anonuuid.FakeUUID(realuuid1)

			So(len(anonuuid.cache), ShouldEqual, 2)
			So(fakeuuid1, ShouldEqual, fakeuuid2)
			So(fakeuuid1, ShouldEqual, fakeuuid5)
			So(fakeuuid3, ShouldEqual, fakeuuid4)
			So(fakeuuid2, ShouldNotEqual, fakeuuid3)
			So(fakeuuid1, ShouldEqual, expected1)
			So(fakeuuid3, ShouldEqual, expected2)
		})
		Convey("--random", func() {
			anonuuid1 := New()
			anonuuid1.Random = true

			fakeuuid11 := anonuuid1.FakeUUID(realuuid1)
			fakeuuid12 := anonuuid1.FakeUUID(realuuid2)

			anonuuid2 := New()
			anonuuid2.Random = true

			fakeuuid21 := anonuuid2.FakeUUID(realuuid1)
			fakeuuid22 := anonuuid2.FakeUUID(realuuid2)

			So(len(anonuuid1.cache), ShouldEqual, 2)
			So(len(anonuuid2.cache), ShouldEqual, 2)
			So(fakeuuid11, ShouldNotEqual, realuuid1)
			So(fakeuuid11, ShouldNotEqual, fakeuuid12)
			So(fakeuuid21, ShouldNotEqual, fakeuuid22)
			So(fakeuuid11, ShouldNotEqual, fakeuuid21)
		})
		Convey("not a valid UUID", func() {
			anonuuid := New()

			output := anonuuid.FakeUUID("hello")
			So(output, ShouldEqual, "invaliduuid")

			output = anonuuid.FakeUUID("hello2")
			So(output, ShouldEqual, "invaliduuid")

			output = anonuuid.FakeUUID("hello")
			So(output, ShouldEqual, "invaliduuid")
		})
		Convey("not a valid UUID with --allow-non-uuid-input", func() {
			anonuuid := New()
			anonuuid.AllowNonUUIDInput = true

			output := anonuuid.FakeUUID("hello")
			So(output, ShouldEqual, "00000000-0000-1000-0000-000000000000")

			output = anonuuid.FakeUUID("hello2")
			So(output, ShouldEqual, "11111111-1111-1111-1111-111111111111")

			output = anonuuid.FakeUUID("hello")
			So(output, ShouldEqual, "00000000-0000-1000-0000-000000000000")
		})
		// FIXME: test cases
		// FIXME: test retry (2 times the same generated uuid)
	})
}

func TestFormatUUID(t *testing.T) {
	Convey("Testing FormatUUID", t, func() {
		out, err := FormatUUID("15573749c89d41dda65516e79bed52e0")
		So(err, ShouldBeNil)
		So(out, ShouldEqual, "15573749-c89d-11dd-a655-16e79bed52e0")

		out, err = FormatUUID("abcdefghijklmnopqrstuvwxyz")
		So(err, ShouldBeNil)
		So(out, ShouldEqual, "abcdefgh-ijkl-1nop-qrst-uvwxyzabcdef")

		out, err = FormatUUID("0123456789")
		So(err, ShouldBeNil)
		So(out, ShouldEqual, "01234567-8901-1345-6789-012345678901")

		out, err = FormatUUID("abcdefghijklmnopqrstuvwxyz0123456789")
		So(err, ShouldBeNil)
		So(out, ShouldEqual, "abcdefgh-ijkl-1nop-qrst-uvwxyz012345")

		out, err = FormatUUID("a")
		So(err, ShouldBeNil)
		So(out, ShouldEqual, "aaaaaaaa-aaaa-1aaa-aaaa-aaaaaaaaaaaa")

		out, err = FormatUUID("@")
		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, "")

		out, err = FormatUUID("")
		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, "")

		out, err = FormatUUID("^")
		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, "")

		/*
			out, err = FormatUUID("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
			So(err, ShouldBeNil)
			So(out, ShouldEqual, "ABCDEFGH-IJKL-1NOP-QRST-UVWXYZABCDEF")
		*/
	})
}

func TestGenerateRandomUUID(t *testing.T) {
	Convey("Testing Anonuuid.GenerateRandomUUID", t, func() {
		out1, err := GenerateRandomUUID(42)
		So(err, ShouldBeNil)
		So(len(out1), ShouldEqual, 36)

		out2, err := GenerateRandomUUID(42)
		So(err, ShouldBeNil)
		So(len(out2), ShouldEqual, 36)
		So(out2, ShouldNotEqual, out1)

		out3, err := GenerateRandomUUID(10)
		So(err, ShouldBeNil)
		So(len(out3), ShouldEqual, 36)
		So(out3, ShouldNotEqual, out2)
		So(out3, ShouldNotEqual, out1)
	})
}

func TestPrefixUUID(t *testing.T) {
	Convey("Testing PrefixUUID", t, func() {
		realuuid := "15573749-c89d-11dd-a655-16e79bed52e0"

		out, err := PrefixUUID("prefix", realuuid)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "prefix15-5737-19c8-9d11-dda65516e79b")

		out, err = PrefixUUID("", realuuid)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, realuuid)

		out, err = PrefixUUID("iamaveryveryveryveryverylongprefix", realuuid)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "iamavery-very-1ery-very-verylongpref")

		out, err = PrefixUUID("@", realuuid)
		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, "")
	})
}

func TestGenerateHexspeakUUID(t *testing.T) {
	Convey("Testing Anonuuid.GenerateHexspeakUUID", t, func() {
		out, err := GenerateHexspeakUUID(0)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "0ff1ce0f-f1ce-1ff1-ce0f-f1ce0ff1ce0f")

		out, err = GenerateHexspeakUUID(0)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "0ff1ce0f-f1ce-1ff1-ce0f-f1ce0ff1ce0f")

		out, err = GenerateHexspeakUUID(1)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "31337313-3731-1373-1337-313373133731")

		out, err = GenerateHexspeakUUID(-1)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "31337313-3731-1373-1337-313373133731")

		// FIXME: i > amount-of-words
	})
}

func TestGenerateLenUUID(t *testing.T) {
	Convey("Testing Anonuuid.GenerateLenUUID", t, func() {
		out, err := GenerateLenUUID(42)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "2a2a2a2a-2a2a-1a2a-2a2a-2a2a2a2a2a2a")

		out, err = GenerateLenUUID(42)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "2a2a2a2a-2a2a-1a2a-2a2a-2a2a2a2a2a2a")

		out, err = GenerateLenUUID(10)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "aaaaaaaa-aaaa-1aaa-aaaa-aaaaaaaaaaaa")

		out, err = GenerateLenUUID(0)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "00000000-0000-1000-0000-000000000000")

		out, err = GenerateLenUUID(100000000000)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "174876e8-0017-1876-e800-174876e80017")

		out, err = GenerateLenUUID(-1)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "3fffffff-3fff-1fff-3fff-ffff3fffffff")

		out, err = GenerateLenUUID(-2)
		So(err, ShouldBeNil)
		So(len(out), ShouldEqual, 36)
		So(out, ShouldEqual, "3ffffffe-3fff-1ffe-3fff-fffe3ffffffe")
	})
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
	// 00000000-0000-1000-0000-000000000000
	// 00000000-0000-1000-0000-000000000000
	// 11111111-1111-1111-1111-111111111111
	// 11111111-1111-1111-1111-111111111111
	// 00000000-0000-1000-0000-000000000000
	// 22222222-2222-1222-2222-222222222222
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

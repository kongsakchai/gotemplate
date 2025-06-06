package logger

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"
)

type A interface {
	String() string
}

type Data struct {
	Value string
}

func (d Data) String() string {
	return d.Value
}

func TestEncodeStyle(t *testing.T) {
	t.Run("should return style when disableColo is false", func(t *testing.T) {
		opt := HandlerOptions{DisableColor: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test styling with color
		encoder.style(buf, txtRed)
		if string(*buf) != txtRed {
			t.Errorf("Expected buffer to contain '%s', got '%s'", txtRed, string(*buf))
		}

		// Test reset
		encoder.reset(buf)
		if string(*buf) != txtRed+txtReset {
			t.Errorf("Expected buffer to contain '%s', got '%s'", txtRed+txtReset, string(*buf))
		}
	})

	t.Run("should return empty buffer when disableColor is true", func(t *testing.T) {
		// Test disable color
		opt := HandlerOptions{DisableColor: true}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test styling with color disabled
		encoder.style(buf, txtRed)
		if string(*buf) != "" {
			t.Errorf("Expected buffer to be empty, got '%s'", string(*buf))
		}

		// Test reset
		encoder.reset(buf)
		if string(*buf) != "" {
			t.Errorf("Expected buffer to be empty after reset, got '%s'", string(*buf))
		}
	})
}

func TestEncodeEmojiLevel(t *testing.T) {
	t.Run("should write emoji ❌ when error level", func(t *testing.T) {
		opt := HandlerOptions{DisableEmoji: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test emoji for error level
		encoder.writeEmojiLevel(buf, slog.LevelError)
		if string(*buf) != "❌ " {
			t.Errorf("Expected buffer to contain '❌ ', got '%s'", string(*buf))
		}
	})

	t.Run("should write emoji ⚠️ when warn level", func(t *testing.T) {
		opt := HandlerOptions{DisableEmoji: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test emoji for warn level
		encoder.writeEmojiLevel(buf, slog.LevelWarn)
		if string(*buf) != "⚠️  " {
			t.Errorf("Expected buffer to contain '⚠️  ', got '%s'", string(*buf))
		}
	})

	t.Run("should write emoji 🌱 when info level", func(t *testing.T) {
		opt := HandlerOptions{DisableEmoji: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test emoji for info level
		encoder.writeEmojiLevel(buf, slog.LevelInfo)
		if string(*buf) != "🌱 " {
			t.Errorf("Expected buffer to contain '🌱 ', got '%s'", string(*buf))
		}
	})

	t.Run("should write emoji 🐛 when debug level", func(t *testing.T) {
		opt := HandlerOptions{DisableEmoji: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test emoji for debug level
		encoder.writeEmojiLevel(buf, slog.LevelDebug)
		if string(*buf) != "🐛 " {
			t.Errorf("Expected buffer to contain '🐛 ', got '%s'", string(*buf))
		}
	})

	t.Run("should not write emoji when disableEmoji is true", func(t *testing.T) {
		opt := HandlerOptions{DisableEmoji: true}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test emoji with disableEmoji
		encoder.writeEmojiLevel(buf, slog.LevelError)
		if string(*buf) != "" {
			t.Errorf("Expected buffer to be empty when disableEmoji is true, got '%s'", string(*buf))
		}
	})
}

func TestEncodeTime(t *testing.T) {
	t.Run("should write time with format when disableTime is false", func(t *testing.T) {
		opt := HandlerOptions{DisableTime: false, TimeFormat: "2006-01-02 15:04:05"}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing time
		testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
		encoder.writeTime(buf, testTime)
		expected := txtGray + "2023-10-01 12:00:00" + txtReset + " "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should not write time when disableTime is true", func(t *testing.T) {
		opt := HandlerOptions{DisableTime: true}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test disable time
		testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
		encoder.writeTime(buf, testTime)
		if string(*buf) != "" {
			t.Errorf("Expected buffer to be empty when time is disabled, got '%s'", string(*buf))
		}
	})
}

func TestEncodeLevel(t *testing.T) {
	t.Run("should write level ERR with read text when error level", func(t *testing.T) {
		opt := HandlerOptions{DisableLevel: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing level with style
		encoder.writeLevel(buf, slog.LevelError)
		expected := txtRed + "ERR " + txtReset
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write level WRN with yellow text when warn level", func(t *testing.T) {
		opt := HandlerOptions{DisableLevel: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing level with style
		encoder.writeLevel(buf, slog.LevelWarn)
		expected := txtYellow + "WRN " + txtReset
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write level INF with green text when info level", func(t *testing.T) {
		opt := HandlerOptions{DisableLevel: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing level with style
		encoder.writeLevel(buf, slog.LevelInfo)
		expected := txtGreen + "INF " + txtReset
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write level DBG with normal text when debug level", func(t *testing.T) {
		opt := HandlerOptions{DisableLevel: false}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing level with style
		encoder.writeLevel(buf, slog.LevelDebug)
		expected := "DBG "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should not write level when disableLevel is true", func(t *testing.T) {
		opt := HandlerOptions{DisableLevel: true}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing level with disableLevel
		encoder.writeLevel(buf, slog.LevelError)
		if string(*buf) != "" {
			t.Errorf("Expected buffer to be empty when disableLevel is true, got '%s'", string(*buf))
		}
	})
}

func TestEncodeMessage(t *testing.T) {
	opt := HandlerOptions{}
	encoder := newEncodeText(opt)
	buf := newBuffer()
	defer buf.Free()

	// Test writing message
	testMessage := "This is a test message"
	encoder.writeMessage(buf, testMessage)
	expected := txtBold + "\"" + testMessage + "\"" + txtReset + " "
	if string(*buf) != expected {
		t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
	}
}

func TestWriteData(t *testing.T) {
	t.Run("should write new line", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing new line
		encoder.writeNewline(buf)
		expected := "\n"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write boolean data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing boolean data
		encoder.writeBool(buf, true)
		expected := "true"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}

		buf.Free()
		encoder.writeBool(buf, false)
		expected = "false"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write integer data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing integer data
		encoder.writeInt(buf, 42)
		expected := "42"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should writeuint data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing uint data
		encoder.writeUint(buf, uint64(42), 10)
		expected := "42"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write float data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing float data
		encoder.writeFloat(buf, 3.14)
		expected := "3.14"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write string data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing string data
		testString := "Hello, World!"
		encoder.writeString(buf, testString)
		expected := "\"" + testString + "\""
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Duration data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Duration data
		testDuration := time.Second * 5
		encoder.writeDuration(buf, testDuration)
		expected := "5000000000"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Time Value data", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Time Value data
		testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
		encoder.writeTimeRFC3339(buf, testTime)
		expected := "2023-10-01T12:00:00.000Z"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})
}

func TestWriteAny(t *testing.T) {
	t.Run("should write bool when reflect.TypeOf is bool", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing bool value
		testBool := true
		val := reflect.ValueOf(testBool)
		encoder.writeAny(buf, val)
		expected := "true"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write int when reflect.TypeOf is int", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing int value
		testInt := 42
		val := reflect.ValueOf(testInt)
		encoder.writeAny(buf, val)
		expected := "42"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write uint when reflect.TypeOf is uint", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing uint value
		testUint := uint(42)
		val := reflect.ValueOf(testUint)
		encoder.writeAny(buf, val)
		expected := "42"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write float when reflect.TypeOf is float", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing float value
		testFloat := 3.14
		val := reflect.ValueOf(testFloat)
		encoder.writeAny(buf, val)
		expected := "3.14"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write string when reflect.TypeOf is string", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing string value
		testString := "Hello, World!"
		val := reflect.ValueOf(testString)
		encoder.writeAny(buf, val)
		expected := "\"" + testString + "\""
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Duration when reflect.TypeOf is time.Duration", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Duration value
		testDuration := time.Second * 5
		val := reflect.ValueOf(testDuration)
		encoder.writeAny(buf, val)
		expected := "5000000000"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Time when reflect.TypeOf is time.Time", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Time value
		testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
		val := reflect.ValueOf(testTime)
		encoder.writeAny(buf, val)
		expected := "2023-10-01 12:00:00 +0000 UTC"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write struct when reflect.TypeOf is struct", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing struct value
		testStruct := struct {
			Name         string
			Age          int
			privateField string // This field should not be exported
		}{
			Name:         "Test Struct",
			Age:          30,
			privateField: "This should not be visible",
		}
		val := reflect.ValueOf(testStruct)
		encoder.writeAny(buf, val)
		expected := "{Name:\"Test Struct\" Age:30 privateField:\"This should not be visible\"}"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write pointer when reflect.TypeOf is slice", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing pointer value
		testSlice := []Data{{Value: "Test Slice 1"}, {Value: "Test Slice 2"}}
		val := reflect.ValueOf(testSlice)
		encoder.writeAny(buf, val)
		expected := "[Test Slice 1 Test Slice 2]"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write map when reflect.TypeOf is map", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing map value
		testMap := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		val := reflect.ValueOf(testMap)
		encoder.writeAny(buf, val)
		expected1 := "\"key1\":\"value1\""
		expected2 := "\"key2\":\"value2\""
		if !strings.Contains(string(*buf), expected1) {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected1, string(*buf))
		}

		if !strings.Contains(string(*buf), expected2) {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected2, string(*buf))
		}
	})

	t.Run("should write interface when reflect.TypeOf is interface", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		var testInterface A = Data{Value: "Test Interface"}

		val := reflect.ValueOf(testInterface)
		encoder.writeAny(buf, val)
		expected := testInterface.String()
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write ptr when reflect.TypeOf is pointer", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing pointer value
		testPtr := &Data{Value: "Test Pointer"}
		val := reflect.ValueOf(testPtr)
		encoder.writeAny(buf, val)
		expected := fmt.Sprintf("%p", testPtr)
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write nil when reflect.TypeOf is pointer but nil", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing nil pointer value
		var nilPtr *Data = nil
		val := reflect.ValueOf(nilPtr)
		encoder.writeAny(buf, val)
		expected := "nil"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write nil when reflect.TypeOf is nil", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing nil value
		var nilValue interface{} = nil
		val := reflect.ValueOf(nilValue)
		encoder.writeAny(buf, val)
		expected := "nil"
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})
}

func TestWriteKey(t *testing.T) {
	t.Run("should write key without group", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing key
		testKey := "testKey"
		encoder.writeKey(buf, []string{}, testKey)
		expected := txtCyan + "testKey" + txtGray + "=" + txtReset
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write key with group", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing key with group
		testKey := "testKey"
		testGroup := []string{"group1", "group2"}
		encoder.writeKey(buf, testGroup, testKey)
		expected := txtCyan + "group1.group2.testKey" + txtGray + "=" + txtReset
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})
}

func TestWriteValue(t *testing.T) {
	t.Run("should write bool when value is bool", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing bool value
		testBool := true
		encoder.writeValue(buf, slog.BoolValue(testBool))
		expected := "true "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write int64 when value is int64", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing int64 value
		testInt64 := int64(42)
		encoder.writeValue(buf, slog.Int64Value(testInt64))
		expected := "42 "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write uint64 when value is uint64", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing uint64 value
		testUint64 := uint64(42)
		encoder.writeValue(buf, slog.Uint64Value(testUint64))
		expected := "42 "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write float64 when value is float64", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing float64 value
		testFloat64 := float64(3.14)
		encoder.writeValue(buf, slog.Float64Value(testFloat64))
		expected := "3.14 "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write string when value is string", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing string value
		testString := "Hello, World!"
		encoder.writeValue(buf, slog.StringValue(testString))
		expected := "\"" + testString + "\" "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Duration when value is time.Duration", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Duration value
		testDuration := time.Second * 5
		encoder.writeValue(buf, slog.DurationValue(testDuration))
		expected := "5000000000 "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write Time when value is time.Time", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing Time value
		testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
		encoder.writeValue(buf, slog.TimeValue(testTime))
		expected := "2023-10-01T12:00:00.000Z "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write nil when value is nil", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing nil value
		var nilValue interface{} = nil
		encoder.writeValue(buf, slog.AnyValue(nilValue))
		expected := "nil "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write nil when value is empty", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing empty value
		var emptyValue interface{}
		encoder.writeValue(buf, slog.AnyValue(emptyValue))
		expected := "nil "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})
}

func TestWriteAttr(t *testing.T) {
	t.Run("should write attribute without group", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing attribute
		testValue := slog.String("key", "value")
		encoder.writeAttr(buf, []string{}, testValue)
		expected := txtCyan + "key" + txtGray + "=" + txtReset + "\"value\" "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write attribute with group", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing attribute with group
		testValue := slog.String("key", "value")
		groupValue := slog.Group("group2", testValue)
		testGroup := []string{"group1"}
		encoder.writeAttr(buf, testGroup, groupValue)
		expected := txtCyan + "group1.group2.key" + txtGray + "=" + txtReset + "\"value\" "
		if string(*buf) != expected {
			t.Errorf("Expected buffer to contain '%s', got '%s'", expected, string(*buf))
		}
	})

	t.Run("should write empaty when empaty slog", func(t *testing.T) {
		opt := HandlerOptions{}
		encoder := newEncodeText(opt)
		buf := newBuffer()
		defer buf.Free()

		// Test writing empty slog
		testValue := slog.Attr{}
		testGroup := []string{"group1"}
		encoder.writeAttr(buf, testGroup, testValue)
		expected := ""
		if string(*buf) != expected {
			t.Errorf("Expected buffer to be empty, got '%s'", string(*buf))
		}
	})
}

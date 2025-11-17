package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/cmd"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/gob"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/housework"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/json"
	"github.com/hippo-an/tiny-go-challenges/netpro_99/encode_decode/protobuf"
)

func main() {
	// 테스트 모드 확인
	if len(os.Args) > 1 && os.Args[1] == "test" {
		runTests()
		return
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTests() {
	fmt.Println("=== 인코딩/디코딩 타입별 테스트 시작 ===")

	// 테스트 데이터 생성
	testChores := []*housework.Chore{
		{Description: "설거지하기", Complete: false},
		{Description: "빨래하기", Complete: true},
		{Description: "청소하기", Complete: false},
		{Description: "요리하기", Complete: true},
		{Description: "쓰레기 버리기", Complete: false},
	}

	fmt.Printf("원본 데이터 (%d개 항목):\n", len(testChores))
	for i, chore := range testChores {
		status := " "
		if chore.Complete {
			status = "X"
		}
		fmt.Printf("  [%s] %d: %s\n", status, i, chore.Description)
	}
	fmt.Println()

	// 1. JSON 테스트
	fmt.Println("--- JSON 인코딩/디코딩 테스트 ---")
	testJSON(testChores)

	// 2. Gob 테스트
	fmt.Println("\n--- Gob 인코딩/디코딩 테스트 ---")
	testGob(testChores)

	// 3. Protobuf 테스트
	fmt.Println("\n--- Protobuf 인코딩/디코딩 테스트 ---")
	testProtobuf(testChores)

	// 4. 크기 비교
	fmt.Println("\n--- 인코딩 크기 비교 ---")
	compareSizes(testChores)

	// 5. 타입 간 호환성 테스트
	fmt.Println("\n--- 타입 간 호환성 테스트 ---")
	testCrossTypeCompatibility(testChores)

	fmt.Println("\n=== 모든 테스트 완료 ===")
}

func testJSON(chores []*housework.Chore) {
	buf := new(bytes.Buffer)

	// 인코딩
	if err := json.Flust(buf, chores); err != nil {
		fmt.Printf("  ❌ JSON 인코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ JSON 인코딩 성공 (크기: %d bytes)\n", buf.Len())

	// 디코딩
	decoded, err := json.Load(buf)
	if err != nil {
		fmt.Printf("  ❌ JSON 디코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ JSON 디코딩 성공 (%d개 항목 복원)\n", len(decoded))

	// 검증
	if verifyChores(chores, decoded) {
		fmt.Println("  ✓ 데이터 무결성 검증 성공")
	} else {
		fmt.Println("  ❌ 데이터 무결성 검증 실패")
	}
}

func testGob(chores []*housework.Chore) {
	buf := new(bytes.Buffer)

	// 인코딩
	if err := gob.Flush(buf, chores); err != nil {
		fmt.Printf("  ❌ Gob 인코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ Gob 인코딩 성공 (크기: %d bytes)\n", buf.Len())

	// 디코딩
	decoded, err := gob.Load(buf)
	if err != nil {
		fmt.Printf("  ❌ Gob 디코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ Gob 디코딩 성공 (%d개 항목 복원)\n", len(decoded))

	// 검증
	if verifyChores(chores, decoded) {
		fmt.Println("  ✓ 데이터 무결성 검증 성공")
	} else {
		fmt.Println("  ❌ 데이터 무결성 검증 실패")
	}
}

func testProtobuf(chores []*housework.Chore) {
	buf := new(bytes.Buffer)

	// 인코딩
	if err := protobuf.Flush(buf, chores); err != nil {
		fmt.Printf("  ❌ Protobuf 인코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ Protobuf 인코딩 성공 (크기: %d bytes)\n", buf.Len())

	// 디코딩
	decoded, err := protobuf.Load(buf)
	if err != nil {
		fmt.Printf("  ❌ Protobuf 디코딩 실패: %v\n", err)
		return
	}
	fmt.Printf("  ✓ Protobuf 디코딩 성공 (%d개 항목 복원)\n", len(decoded))

	// 검증
	if verifyChores(chores, decoded) {
		fmt.Println("  ✓ 데이터 무결성 검증 성공")
	} else {
		fmt.Println("  ❌ 데이터 무결성 검증 실패")
	}
}

func compareSizes(chores []*housework.Chore) {
	// JSON 크기
	jsonBuf := new(bytes.Buffer)
	json.Flust(jsonBuf, chores)
	jsonSize := jsonBuf.Len()

	// Gob 크기
	gobBuf := new(bytes.Buffer)
	gob.Flush(gobBuf, chores)
	gobSize := gobBuf.Len()

	// Protobuf 크기
	protoBuf := new(bytes.Buffer)
	protobuf.Flush(protoBuf, chores)
	protoSize := protoBuf.Len()

	fmt.Printf("  JSON:     %5d bytes (100.0%%)\n", jsonSize)
	fmt.Printf("  Gob:      %5d bytes (%5.1f%%)\n", gobSize, float64(gobSize)/float64(jsonSize)*100)
	fmt.Printf("  Protobuf: %5d bytes (%5.1f%%)\n", protoSize, float64(protoSize)/float64(jsonSize)*100)

	// 가장 작은 크기 찾기
	minSize := jsonSize
	minType := "JSON"
	if gobSize < minSize {
		minSize = gobSize
		minType = "Gob"
	}
	if protoSize < minSize {
		minSize = protoSize
		minType = "Protobuf"
	}
	fmt.Printf("\n  가장 효율적인 인코딩: %s (%d bytes)\n", minType, minSize)
}

func testCrossTypeCompatibility(chores []*housework.Chore) {
	fmt.Println("  JSON으로 인코딩 후 Gob으로 디코딩 시도...")
	jsonBuf := new(bytes.Buffer)
	json.Flust(jsonBuf, chores)
	if _, err := gob.Load(jsonBuf); err != nil {
		fmt.Printf("    ✓ 예상대로 실패: %v\n", err)
	} else {
		fmt.Println("    ❌ 예상과 다르게 성공 (호환되지 않아야 함)")
	}

	fmt.Println("\n  Gob으로 인코딩 후 Protobuf으로 디코딩 시도...")
	gobBuf := new(bytes.Buffer)
	gob.Flush(gobBuf, chores)
	if _, err := protobuf.Load(gobBuf); err != nil {
		fmt.Printf("    ✓ 예상대로 실패: %v\n", err)
	} else {
		fmt.Println("    ❌ 예상과 다르게 성공 (호환되지 않아야 함)")
	}

	fmt.Println("\n  Protobuf으로 인코딩 후 JSON으로 디코딩 시도...")
	protoBuf := new(bytes.Buffer)
	protobuf.Flush(protoBuf, chores)
	if _, err := json.Load(protoBuf); err != nil {
		fmt.Printf("    ✓ 예상대로 실패: %v\n", err)
	} else {
		fmt.Println("    ❌ 예상과 다르게 성공 (호환되지 않아야 함)")
	}
}

func verifyChores(original, decoded []*housework.Chore) bool {
	if len(original) != len(decoded) {
		return false
	}

	for i := range original {
		if original[i].Description != decoded[i].Description {
			return false
		}
		if original[i].Complete != decoded[i].Complete {
			return false
		}
	}

	return true
}

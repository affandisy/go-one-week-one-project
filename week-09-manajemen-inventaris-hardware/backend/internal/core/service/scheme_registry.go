package service

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// SchemaValidator bertugas memvalidasi atribut spesifikasi dinamis
type SchemaValidator struct {
	compilers map[string]*jsonschema.Schema
}

// NewSchemaValidator menginisialisasi dan memuat semua skema JSON ke dalam memori (RAM)
func NewSchemaValidator() (*SchemaValidator, error) {
	sv := &SchemaValidator{
		compilers: make(map[string]*jsonschema.Schema),
	}

	// 1. Skema CPU (Contoh ideal untuk Ryzen 5600 atau setaranya)
	cpuSchema := `{
		"type": "object",
		"required": ["socket", "cores", "threads", "tdp_w"],
		"properties": {
			"socket": { "type": "string" },
			"cores": { "type": "integer", "minimum": 1 },
			"threads": { "type": "integer", "minimum": 1 },
			"tdp_w": { "type": "integer", "minimum": 5 },
			"integrated_graphics": { "type": "boolean" }
		},
		"additionalProperties": false
	}`

	// 2. Skema GPU (Contoh ideal untuk menampung data RTX 4060)
	gpuSchema := `{
		"type": "object",
		"required": ["pcie_generation", "vram_gb", "power_draw_w", "recommended_psu_w"],
		"properties": {
			"pcie_generation": { "type": "string" },
			"vram_gb": { "type": "integer", "minimum": 1 },
			"power_draw_w": { "type": "integer", "minimum": 10 },
			"recommended_psu_w": { "type": "integer", "minimum": 100 }
		},
		"additionalProperties": false
	}`

	// 3. Skema Perangkat Seluler / Daily Driver (Xiaomi, Samsung, dll.)
	mobileSchema := `{
		"type": "object",
		"required": ["os", "ram_gb", "storage_gb", "battery_mah"],
		"properties": {
			"os": { "type": "string", "enum": ["Android", "iOS", "HarmonyOS"] },
			"ram_gb": { "type": "integer", "minimum": 1 },
			"storage_gb": { "type": "integer", "minimum": 8 },
			"battery_mah": { "type": "integer", "minimum": 1000 },
			"fast_charging_w": { "type": "integer" }
		},
		"additionalProperties": false
	}`

	// Kompilasi skema saat aplikasi berjalan (startup) untuk performa maksimal
	if err := sv.compileAndAdd("CPU", cpuSchema); err != nil {
		return nil, err
	}
	if err := sv.compileAndAdd("GPU", gpuSchema); err != nil {
		return nil, err
	}
	if err := sv.compileAndAdd("Mobile", mobileSchema); err != nil {
		return nil, err
	}

	return sv, nil
}

func (sv *SchemaValidator) compileAndAdd(category, schemaString string) error {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(category+".json", strings.NewReader(schemaString)); err != nil {
		return err
	}
	schema, err := compiler.Compile(category + ".json")
	if err != nil {
		return err
	}
	sv.compilers[category] = schema
	return nil
}

// Validate mengeksekusi pemeriksaan payload spesifikasi terhadap kategori yang sesuai
func (sv *SchemaValidator) Validate(category string, specs map[string]interface{}) error {
	schema, exists := sv.compilers[category]
	if !exists {
		return fmt.Errorf("skema validasi untuk kategori '%s' tidak ditemukan", category)
	}

	if err := schema.Validate(specs); err != nil {
		// Mengembalikan pesan error validasi JSON Schema yang spesifik
		return fmt.Errorf("spesifikasi tidak valid: %v", err)
	}
	return nil
}

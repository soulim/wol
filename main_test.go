package main

import (
	"testing"
)

func Test_wol(t *testing.T) {
	type args struct {
		macAddress string
	}
	// Defining the columns of the table
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Scenarios
		{
			name: "Valid MAC Address",
			args: args{
				macAddress: "00:1B:44:11:3A:B7", // random valid mac
			},
			wantErr: false,
		},
		{
			name: "Invalid MAC Address",
			args: args{
				macAddress: "IN:VA:LI:D0:0M:AC", // an invalid case
			},
			wantErr: true,
		},
	}

	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wol(tt.args.macAddress)
			if err != nil && !tt.wantErr {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

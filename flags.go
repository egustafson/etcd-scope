package main

// Flags
//
var (
	verboseFlag  = false
	configFlag   = ""
	endpointFlag = DefaultEndpoint
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false,
		"enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&configFlag, "config", "c", "",
		"config file (default is $HOME/.config/etcd-scope.yml)")
	rootCmd.PersistentFlags().StringVarP(&endpointFlag, "endpoint", "e", DefaultEndpoint,
		"etcd endpoint (default: localhost:2379)")
}

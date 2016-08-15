package resources


const (
  CPU_USAGE_PERCENTAGE = "top -b -n2 -d 0.75 | grep \"Cpu(s)\" | tail -n 1 | awk '{print $2 + $4}'"
  MEMORY_USAGE_MB = "free -m | grep \"Mem:\" | awk '{print $3}'"
  MEMORY_USAGE_PERCENT = "free -m | grep \"Mem:\" | awk '{print ($3 / $2) * 100}'"
  SWAP_USAGE_MB = "free -m | grep \"Swap:\" | awk '{print $3}'"
  // SWAP_USAGE_PERCENT = "free -m | grep \"Swap:\" | awk '{print ($3 / $2) * 100}'"
)

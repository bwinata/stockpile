package resources


const (
  INDEX_CPU         = 0
  INDEX_MEM_USED    = 1
  INDEX_MEM_TOTAL   = 2

  CPU_USAGE_PERCENTAGE = "top -b -n2 -d 0.5 | grep \"Cpu(s)\" | tail -n 1 | awk '{print $2 + $4}'"
  MEMORY_USAGE_PERCENTAGE = ""
  MEMORY_USAGE_MB = "free -m | grep \"Mem:\" | awk '{print $3}'"
  MEMORY_TOTAL_MB = "free -m | grep \"Mem:\" | awk '{print $2}'"
)

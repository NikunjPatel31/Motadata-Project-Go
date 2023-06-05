package util

const Status = "status"

const Message = "message"

const SystemMemoryInstalledBytes = "system.memory.installed.bytes"

const SystemMemoryUsedBytes = "system.memory.used.bytes"

const SystemMemoryUsedPercentage = "system.memory.used.percentage"

const SystemMemoryFreeBytes = "system.memory.free.bytes"

const SystemMemoryFreePercentage = "system.memory.free.percentage"

const SystemMemoryAvailableBytes = "system.memory.available.bytes"

const SystemMemorySwapBytes = "system.memory.swap.bytes"

const SystemMemoryTotalBytes = "system.memory.total.bytes"

const SystemProcessPid = "system.process.pid"

const SystemProcessCPU = "system.process.cpu"

const SystemProcessMemory = "system.process.memory"

const SystemProcessUser = "system.process.user"

const SystemProcessCommand = "system.process.command"

const cmd = "nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 \" \" $7 \" \" $5 \" \" $14}'"

const SystemCPUPercentage = "system.cpu.percentage"

const SystemCPUCore = "system.cpu.core"

const SystemCPUUserPercentage = "system.cpu.user.percentage"

const SystemCPUIdlePercentage = "system.cpu.idle.percentage"

const SystemDisk = "system.disk"

const SystemCPU = "system.cpu"

const SystemProcess = "system.process"

const SystemName = "system.name"

const SystemOSName = "system.os.name"

const SystemOSVersion = "system.os.version"

const SystemUptime = "system.uptime"

const SystemThreads = "system.threads"

const SystemContextSwtiches = "system.context.switches"

const SystemRunningProcesses = "system.running.processes"

const SystemBlockedProcesses = "system.blocked.processes"

const NewLineSeparator = "\n"

const SpaceSeparator = " "

const Input = "input"

const Result = "result"

const Fail = "fail"

const Success = "success"

const CPU = "cpu"

const Memory = "memory"

const Process = "process"

const System = "system"

const Disk = "disk"

-- setup thread
local thread_num = 0
function setup(thread)
	thread_num = thread_num + 1
end

-- Get args from command line
-- args[0]: echo size
local echo_size = 0
function init(args)
	echo_size = args[1]
	-- prepare request
	wrk.method = "POST"
	b = ""
	for i = 1, echo_size
	do
		b = b .. "1"
	end
	wrk.body = b
	wrk.headers["Content-Type"] = "application/json"
end

-- report data
result = {}
function done(summary, latency, requests)
	io.write("----------------\n")
	io.write("Benchmark_Result\n")
	-- -- Concurrency
	-- result["Concurrency"] = thread_num
	-- io.write(string.format("Concurrency,%d\n", thread_num))
	
	-- QPS
	local throughput = summary.requests/(summary.duration/(10^6))
	result["QPS"] = throughput
	io.write(string.format("QPS,%d\n", throughput))

	-- Latency, TP99, TP999
	for _, p in pairs({ 99, 99.9 }) do
		local n = latency:percentile(p)
		result[string.format("TP%g", p)] = n
		io.write(string.format("TP%g,%d\n", p, n))
	end
	
	-- Error
	-- local read_error = summary.errors.read
	-- local connect_error = summary.errors.connect
	-- local write_error = summary.errors.write
	-- local status_error = summary.errors.status
	-- local timeout_error = summary.errors.timeout
	local total_error = 0
	for i, value in pairs(summary.errors) do
		total_error = total_error + value
	end
	io.write(string.format("request,%d\n", summary.requests))
	io.write(string.format("failed,%d\n", total_error))
	io.write("----------------\n\n")
end
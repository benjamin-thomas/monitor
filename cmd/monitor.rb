#!/usr/bin/env ruby

require 'colored'
require 'logger'
require 'etc'

NUM_CPUS=Etc.nprocessors

MAX_AVG = 1 * NUM_CPUS
MAX_MEM_PERCENT = 80

LOG = Logger.new(STDOUT)

f = File.read("/proc/loadavg")
LoadAvg = Struct.new(:one_minute, :five_minutes, :fifteen_minutes)

class MemStats
  # $ cat /proc/meminfo | grep ^Mem
  # MemTotal:        4052056 kB
  # MemFree:          249816 kB
  # MemAvailable:    1475728 kB

  def initialize
    @mem = `cat /proc/meminfo | grep ^Mem`
  end

  def total
    # man free: total  Total installed memory (MemTotal and SwapTotal in /proc/meminfo)
    @total ||= extract('MemTotal:')
  end

  def free
    # man free: free Unused memory (MemFree and SwapFree in /proc/meminfo)
    # This value should be mostly ignored, "available" is more correct to determine the free memory available
    @free ||= extract('MemFree:')
  end

  def available
    # man free: Estimation  of how much memory is available for starting new applications, without swapping
    @available ||= extract('MemAvailable:')
  end

  def used
    # Matches htop's output
    @used ||= total - available
  end

  def percent_used
    # Matches top's output (press 'm' several times)
    @percent_used ||= used / Float(total) * 100
  end

  def inspect
    pretty = ->(num) { "#{num} (%.2f GB)" % (num / 1024.0**2) }
    "#{super.inspect.split[0]} total=#{pretty.call(total)}, free=#{pretty.call(free)}, available=#{pretty.call(available)}, used=#{pretty.call(used)}, percent_used=#{percent_used.round(2)}%"
  end

  private

  def extract(label)
    str_int = @mem.split("\n").find { |line| line.start_with?(label) }.split[1]
    Integer(str_int)
  end
end

avgs_floats = f.split[0,3].map { |str| Float(str) }

avg = LoadAvg.new(*avgs_floats)

if avg.fifteen_minutes > MAX_AVG
  LOG.info(
    "Avg is too high!! (avg=#{avg}, max=#{MAX_AVG})".red
  )
else
  LOG.info(
    "Avg is OK (avg=#{avg}, max=#{MAX_AVG})".green
  )
end

ms = MemStats.new

if ms.percent_used > MAX_MEM_PERCENT
  LOG.info("Memory usage is too high!! (#{ms.inspect})".red)
else
  LOG.info("Memory usage is OK (#{ms.inspect})".green)
end

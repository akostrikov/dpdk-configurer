https://github.com/InfraSIM/infrasim-compute/wiki/Performance-Tuning-Tips
https://github.com/snabbco/snabb/blob/master/src/doc/performance-tuning.md

# Cpu tasks to care
CPU pinning via tasket and numactl work also within Docker Containers running in privileged mode. It is important to note that the Container can use all CPU cores, including the ones specifically excluded by the kernel option isolcpus:

Watch for TLB shootdowns

# apparmor/selinux

# add autoloading 
echo uio >> /etc/modules
echo uio_pci_generic >> /etc/modules
echo igb_uio >> /etc/modules
echo rte_kni >> /etc/modules

# kernel params
cat /sys/devices/system/cpu/possible
GRUB_CMDLINE_LINUX_DEFAULT="isolcpus=2-19 default_hugepagesz=1GB hugepagesz=1G hugepages=32 intel_iommu=off intel_idle.max_cstate=0"


# disable irqbalance
systemctl disable irqbalance
systemctl stop irqbalance
sed -i s/'ENABLED="1"'/'ENABLED="0"'/ /etc/default/irqbalance

# hugepages count (numa?)

# check kernel parameters
cat /sys/devices/system/cpu/isolated
taskset -cp 1
grep "Hugepagesize:" /proc/meminfo

# show cpus by numa
# show interfaces by numa
commands to bind to driver.

# cpu perfomance mode
```
for CPUFREQ in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor;

do [ -f $CPUFREQ ] || continue;
echo -n performance > $CPUFREQ;
done
grep -E '^model name|^cpu MHz' /proc/cpuinfo
```

# numa bind
```
option: numactl --physcpubind=xxxx --localalloc qemu-application Filtering all core_id/socket_id before you bind.

cores = [0, 1, 2, 3, 4, 8, 9, 10, 11, 12] sockets = [0, 1] Socket 0 Socket 1 -------- -------- Core 0 [0, 20] [10, 30] Core 1 [1, 21] [11, 31] Core 2 [2, 22] [12, 32] Core 3 [3, 23] [13, 33] Core 4 [4, 24] [14, 34] Core 8 [5, 25] [15, 35] Core 9 [6, 26] [16, 36] Core 10 [7, 27] [17, 37] Core 11 [8, 28] [18, 38] Core 12 [9, 29] [19, 39] http://www.glennklockwood.com/hpc-howtos/process-affinity.html
```


# compute mask
# pass driver

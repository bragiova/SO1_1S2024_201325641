#include <linux/build-salt.h>
#include <linux/module.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(.gnu.linkonce.this_module) = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif

static const struct modversion_info ____versions[]
__used __section(__versions) = {
	{ 0xa29da427, "module_layout" },
	{ 0x23a1de1b, "seq_read" },
	{ 0x5bbb8213, "remove_proc_entry" },
	{ 0xc5850110, "printk" },
	{ 0x635a6874, "proc_create" },
	{ 0xeae3dfd6, "__const_udelay" },
	{ 0xb1ddf995, "jiffies_64_to_clock_t" },
	{ 0x17de3d5, "nr_cpu_ids" },
	{ 0xc5e4a5d1, "cpumask_next" },
	{ 0x9e683f75, "__cpu_possible_mask" },
	{ 0x1bb467bc, "init_task" },
	{ 0xfc255b09, "seq_printf" },
	{ 0x1234e483, "get_cpu_iowait_time_us" },
	{ 0x75d0deb9, "nsecs_to_jiffies64" },
	{ 0x7b9793a2, "get_cpu_idle_time_us" },
	{ 0xb19a5453, "__per_cpu_offset" },
	{ 0xb58aeaab, "kernel_cpustat" },
	{ 0x5a5a2271, "__cpu_online_mask" },
	{ 0x8f0cb4bb, "single_open" },
	{ 0xbdfb6dbb, "__fentry__" },
};

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "7076C0BBC31714BF3FDB1B0");

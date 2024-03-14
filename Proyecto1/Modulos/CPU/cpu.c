// Info de los modulos
#include <linux/module.h>
// Info del kernel en tiempo real
#include <linux/sched/signal.h>
#include <linux/kernel.h>
#include <linux/sched.h>

// Headers para modulos
#include <linux/init.h>
// Header necesario para proc_fs
#include <linux/proc_fs.h>
// Para dar acceso al usuario
#include <asm/uaccess.h>
// Para manejar el directorio /proc
#include <linux/seq_file.h>
// Para get_mm_rss
#include <linux/mm.h>
#include <linux/fs.h>

#include <linux/vmstat.h>
#include <linux/kernel_stat.h>
#include <linux/delay.h>
#include <linux/version.h>

#include <linux/tick.h>

struct task_struct *task; // Estructura que almacena info del cpu

// Almacena los procesos
struct list_head *lstProcess;
// Estructura que almacena info de los procesos hijos
struct task_struct *taskChild;
unsigned long rss;

unsigned cpu_stat = 0;

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo de CPU SO1");
MODULE_AUTHOR("BrayanRivas");

//Metodos obtenidos de stat.c
static u64 get_idle_time(int cpu_st)
{
	u64 idle, idle_time = -1ULL;

	if (cpu_online(cpu_st))
		idle_time = get_cpu_idle_time_us(cpu_st, NULL);

	if (idle_time == -1ULL)
		idle = kcpustat_cpu(cpu_st).cpustat[CPUTIME_IDLE];
	else
		idle = nsecs_to_jiffies64(idle_time);

	return idle;
}

//Metodos obtenidos de stat.c
static u64 get_iowait_time(int cpu_st)
{
	u64 iowait, iowait_time = -1ULL;

	if (cpu_online(cpu_st))
		iowait_time = get_cpu_iowait_time_us(cpu_st, NULL);

	if (iowait_time == -1ULL)
		iowait = kcpustat_cpu(cpu_st).cpustat[CPUTIME_IOWAIT];
	else
		iowait = nsecs_to_jiffies64(iowait_time);

	return iowait;
}

static int cpu_stat_show(struct seq_file *m, void *v){
    int i;
    u64 Total = 0;
	unsigned long long total_jiffies_1, total_jiffies_2;
	unsigned long long work_jiffies_1, work_jiffies_2;
	unsigned int work_over_period, total_over_period;
	unsigned int usage;

    u64 user, nice, system, idle, iowait, irq, softirq, steal;
 	u64 guest, guest_nice;
    
    user = nice = system = idle = iowait = irq = softirq = steal = 0;
	guest = guest_nice = 0;
    //recorremos todos los cpus y sumamos sus procesos
    for_each_possible_cpu(i) {
		user += kcpustat_cpu(i).cpustat[CPUTIME_USER];
		nice += kcpustat_cpu(i).cpustat[CPUTIME_NICE];
		system += kcpustat_cpu(i).cpustat[CPUTIME_SYSTEM];
		idle += get_idle_time(i);
		iowait += get_iowait_time(i);
		irq += kcpustat_cpu(i).cpustat[CPUTIME_IRQ];
		softirq += kcpustat_cpu(i).cpustat[CPUTIME_SOFTIRQ];
		steal += kcpustat_cpu(i).cpustat[CPUTIME_STEAL];
		guest += kcpustat_cpu(i).cpustat[CPUTIME_GUEST];
		guest_nice += kcpustat_cpu(i).cpustat[CPUTIME_GUEST_NICE];
	}
    Total = user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice;
	work_jiffies_1 = jiffies_64_to_clock_t(user+nice+system);
	total_jiffies_1 = jiffies_64_to_clock_t(Total);
    //delay para obtener otra medicion
	mdelay(200);
	user = nice = system = idle = iowait = irq = softirq = steal = 0;
	guest = guest_nice = 0;
    for_each_possible_cpu(i) {
		user += kcpustat_cpu(i).cpustat[CPUTIME_USER];
		nice += kcpustat_cpu(i).cpustat[CPUTIME_NICE];
		system += kcpustat_cpu(i).cpustat[CPUTIME_SYSTEM];
		idle += get_idle_time(i);
		iowait += get_iowait_time(i);
		irq += kcpustat_cpu(i).cpustat[CPUTIME_IRQ];
		softirq += kcpustat_cpu(i).cpustat[CPUTIME_SOFTIRQ];
		steal += kcpustat_cpu(i).cpustat[CPUTIME_STEAL];
		guest += kcpustat_cpu(i).cpustat[CPUTIME_GUEST];
		guest_nice += kcpustat_cpu(i).cpustat[CPUTIME_GUEST_NICE];
	}
    Total = user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice;
	work_jiffies_2 = jiffies_64_to_clock_t(user+nice+system);
	total_jiffies_2 = jiffies_64_to_clock_t(Total);

	work_over_period = work_jiffies_2-work_jiffies_1;
	total_over_period = total_jiffies_2-total_jiffies_1;
	usage = work_over_period*100/total_over_period;

	
	seq_printf(m, "\"cpu\":%d\n",usage);
    return 0;
}

static int escribir_archivo(struct seq_file *archivo, void *v) {
    int running = 0;
    int sleeping = 0;
    int zombie = 0;
    int stopped = 0;

    seq_printf(archivo, "{\n\"processes\":[\n");
    int b = 0;
    for_each_process(task)
    {
        if (task->mm)
        {
            rss = get_mm_rss(task->mm) << PAGE_SHIFT;
        }
        else
        {
            rss = 0;
        }
        if (b == 0)
        {
            seq_printf(archivo, "{");
            b = 1;
        }
        else
        {
            seq_printf(archivo, ",{");
        }
        seq_printf(archivo, "\"pid\":%d,\n", task->pid);
        seq_printf(archivo, "\"name\":\"%s\",\n", task->comm);
        seq_printf(archivo, "\"user\": %d,\n", task->cred->uid);
        seq_printf(archivo, "\"state\":%ld,\n", task->state);
        int porcentaje = (((rss / (1024 * 1024))) * 100) / (15685);
        seq_printf(archivo, "\"ram\":%d,\n", porcentaje);

        seq_printf(archivo, "\"child\":[\n");
        int a = 0;
        list_for_each(lstProcess, &(task->children))
        {
            taskChild = list_entry(lstProcess, struct task_struct, sibling);
            if (a != 0)
            {
                seq_printf(archivo, ",{");
                seq_printf(archivo, "\"pid\":%d,\n", taskChild->pid);
                seq_printf(archivo, "\"name\":\"%s\",\n", taskChild->comm);
                seq_printf(archivo, "\"state\":%ld\n", taskChild->state);
                seq_printf(archivo, "}\n");
            }
            else
            {
                seq_printf(archivo, "{");
                seq_printf(archivo, "\"pid\":%d,\n", taskChild->pid);
                seq_printf(archivo, "\"name\":\"%s\",\n", taskChild->comm);
                seq_printf(archivo, "\"state\":%ld\n", taskChild->state);
                seq_printf(archivo, "}\n");
                a = 1;
            }
        }
        a = 0;
        seq_printf(archivo, "\n]");

        if (task->state == 0)
        {
            running += 1;
        }
        else if (task->state == 1)
        {
            sleeping += 1;
        }
        else if (task->state == 4)
        {
            zombie += 1;
        }
        else
        {
            stopped += 1;
        }
        seq_printf(archivo, "}\n");
    }
    b = 0;
    seq_printf(archivo, "],\n");
    seq_printf(archivo, "\"running\":%d,\n", running);
    seq_printf(archivo, "\"sleeping\":%d,\n", sleeping);
    seq_printf(archivo, "\"zombie\":%d,\n", zombie);
    seq_printf(archivo, "\"stopped\":%d,\n", stopped);
    seq_printf(archivo, "\"total\":%d,\n", running + sleeping + zombie + stopped);
    cpu_stat_show(archivo,v);
    seq_printf(archivo, "}\n");

    // for_each_process(cpu) {
    //     seq_printf(archivo, "PID%d", cpu->pid);
    //     seq_printf(archivo, ",");
    //     seq_printf(archivo, "%s", cpu->comm);
    //     seq_printf(archivo, ",");
    //     seq_printf(archivo, "%lu", cpu->state);
    //     seq_printf(archivo, ",");

    //     if (cpu->mm) {
    //         rss = get_mm_rss(cpu->mm) << PAGE_SHIFT;
    //         seq_printf(archivo, "%lu", rss);
    //     } else {
    //         seq_printf(archivo, "%s", "");
    //     }
    //     seq_printf(archivo, ",");

    //     seq_printf(archivo, "%d", cpu->cred->user->uid);
    //     seq_printf(archivo, ",");

    //     list_for_each(lstProcess, &(cpu->children)) {
    //         child = list_entry(lstProcess, struct task_struct, sibling);
    //         seq_printf(archivo, "Child:%d", child->pid);
    //         seq_printf(archivo, ".");
    //         seq_printf(archivo, "%s", child->comm);
    //         seq_printf(archivo, ".");
    //         seq_printf(archivo, "%d", child->state);
    //         seq_printf(archivo, ".");

    //          if (child->mm) {
    //             rss = get_mm_rss(child->mm) << PAGE_SHIFT;
    //             seq_printf(archivo, "%lu", rss);
    //         } else {
    //             seq_printf(archivo, "%s", "");
    //         }
    //         seq_printf(archivo, ".");

    //         seq_printf(archivo, "%d", child->cred->user->uid);
    //     }
    // }

    return 0;
}

//Funcion que se ejecutara cada vez que se lea el archivo con el comando CAT
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

//Si el kernel es 5.6 o mayor se usa la estructura proc_ops
#ifdef HAVE_PROC_OPS
static struct proc_ops operaciones =
{
    .proc_open = al_abrir,
    .proc_read = seq_read
};
#else
static const struct file_operations operaciones = {
    .owner = THIS_MODULE,
    .open = al_abrir,
    .read = seq_read
    // .llseek = seq_lseek,
    // .release = single_release,
};
#endif

//Funcion a ejecuta al insertar el modulo en el kernel con insmod
static int _insert(void)
{
    proc_create("cpu_so1_1s2024", 0, NULL, &operaciones);
    printk(KERN_INFO "Brayan Giovanny Rivas Estrada\n");
    return 0;
}

//Funcion a ejecuta al remover el modulo del kernel con rmmod
static void _remove(void)
{
    remove_proc_entry("cpu_so1_1s2024", NULL);
    printk(KERN_INFO "Primer Semestre 2024\n");
}

module_init(_insert);
module_exit(_remove);

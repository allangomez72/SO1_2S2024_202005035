#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tu Nombre");
MODULE_DESCRIPTION("Modulo para leer informacion de memoria y procesos");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo"

static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo si;
    struct task_struct *task;

    si_meminfo(&si);

    seq_printf(m, "Total RAM: %lu KB\n", si.totalram * 4);
    seq_printf(m, "Free RAM: %lu KB\n", si.freeram * 4);
    seq_printf(m, "Number of processes: %d\n", num_online_cpus());

    seq_printf(m, "\nProcess and Subprocess Information:\n");

    for_each_process(task) {
        if (task->parent) {
            seq_printf(m, "Parent Process %s [PID: %d]\n", task->parent->comm, task->parent->pid);
            seq_printf(m, " Child Process %s [PID: %d]\n", task->comm, task->pid);
        }
    }

    return 0;
}

static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};

static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "sysinfo module loaded\n");
    return 0;
}

static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo module unloaded\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);

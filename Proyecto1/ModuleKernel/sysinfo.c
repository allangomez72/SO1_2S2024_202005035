#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/string.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched.h>
#include <linux/timer.h>
#include <linux/jiffies.h>
#include <linux/uaccess.h>
#include <linux/tty.h>
#include <linux/sched/signal.h>
#include <linux/fs.h>
#include <linux/slab.h>
#include <linux/sched/mm.h>
#include <linux/binfmts.h>
#include <linux/timekeeping.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Allan_Gomez");
MODULE_DESCRIPTION("Modulo para leer informacion de memoria y CPU en JSON");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo_202005035"
#define MAX_CMDLINE_LENGTH 256
#define CONTAINER_ID_LENGTH 64

// Función para obtener la línea de comandos de un proceso y retornar un apuntador a la cadena
static char *get_process_cmdline(struct task_struct *task) {
    struct mm_struct *mm;
    char *cmdline, *p;
    unsigned long arg_start, arg_end, env_start;
    int i, len;

    // Reservamos memoria para la línea de comandos
    cmdline = kmalloc(MAX_CMDLINE_LENGTH, GFP_KERNEL);
    if (!cmdline)
        return NULL;

    // Obtenemos la estructura mm_struct del proceso para acceder a la memoria
    mm = get_task_mm(task);
    if (!mm) {
        kfree(cmdline);
        return NULL;
    }

    // Bloqueamos la lectura de mm_struct para acceso seguro
    down_read(&mm->mmap_lock);
    arg_start = mm->arg_start;
    arg_end = mm->arg_end;
    env_start = mm->env_start;
    up_read(&mm->mmap_lock);

    // Calculamos la longitud de la línea de comandos y aseguramos que no exceda el límite
    len = arg_end - arg_start;
    if (len > MAX_CMDLINE_LENGTH - 1)
        len = MAX_CMDLINE_LENGTH - 1;

    // Leemos la línea de comandos de la memoria del proceso
    if (access_process_vm(task, arg_start, cmdline, len, 0) != len) {
        mmput(mm);
        kfree(cmdline);
        return NULL;
    }

    // Añadimos un carácter nulo al final de la línea de comandos y reemplazamos caracteres nulos por espacios
    cmdline[len] = '\0';
    p = cmdline;
    for (i = 0; i < len; i++)
        if (p[i] == '\0')
            p[i] = ' ';

    mmput(mm);
    return cmdline;
}

// Función para mostrar la información en formato JSON en el archivo /proc
static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo si;
    struct task_struct *task;
    unsigned long total_jiffies = jiffies;
    int first_process = 1;

    // Obtenemos la información del sistema, incluyendo memoria
    si_meminfo(&si);

    // Imprimimos la información en formato JSON
    seq_printf(m, "{\n");
    seq_printf(m, "\"Total RAM\": %lu,\n", si.totalram * 4);  // Total de RAM en KB
    seq_printf(m, "\"Free RAM\": %lu,\n", si.freeram * 4);    // RAM libre en KB
    seq_printf(m, "\"Used Ram\": %lu,\n", (si.totalram - si.freeram) * 4);  // RAM usada en KB
    seq_printf(m, "\"Processes\": [\n");

    // Iteramos sobre todos los procesos
    for_each_process(task) {
        // Solo procesamos procesos cuyo nombre es "containerd-shim"
        if (strcmp(task->comm, "containerd-shim") == 0) {
            struct list_head *list;
            struct task_struct *child;
            unsigned long vsz = 0;
            unsigned long rss = 0;
            unsigned long totalram = si.totalram * 4;
            unsigned long mem_usage = 0;
            unsigned long cpu_usage = 0;
            char *cmdline = NULL;

            // Obtenemos el uso de memoria virtual y residente
            if (task->mm) {
                vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
                rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);
                mem_usage = (rss * 10000) / totalram;  // Calculamos el uso de memoria en porcentaje
            }

            // Calculamos el uso de CPU en porcentaje
            unsigned long total_time = task->utime + task->stime;
            cpu_usage = (total_time * 10000) / total_jiffies;
            
            // Recorremos la lista de procesos hijos de containerd-shim
            list_for_each(list, &task->children) {
                child = list_entry(list, struct task_struct, sibling);

                // Solo mostramos procesos que tienen nombre "python" o "high_ram.py"
                if (strcmp(child->comm, "python") == 0 ||  strcmp(child->comm, "sleep") == 0) {
                    cmdline = get_process_cmdline(task);

                    // Imprimimos información del proceso hijo en formato JSON
                    if (!first_process) {
                        seq_printf(m, ",\n");
                    } else {
                        first_process = 0;
                    }

                    seq_printf(m, "  {\n");
                    seq_printf(m, "    \"PID\": %d,\n", child->pid);
                    seq_printf(m, "    \"Name\": \"%s\",\n", child->comm);
                    seq_printf(m, "    \"Cmdline\": \"%s\",\n", cmdline ? cmdline : "N/A");
                    seq_printf(m, "    \"VSZ\": %lu,\n", vsz);  // Tamaño del proceso en KB
                    seq_printf(m, "    \"RSS\": %lu,\n", rss);  // Uso de memoria residente en KB
                    seq_printf(m, "    \"MemoryUsage\": %lu.%02lu,\n", mem_usage / 100, mem_usage % 100);  // Uso de memoria en porcentaje
                    seq_printf(m, "    \"CPUUsage\": %lu.%02lu\n", cpu_usage / 100, cpu_usage % 100);  // Uso de CPU en porcentaje
                    seq_printf(m, " }");

                    // Liberamos la memoria de la línea de comandos
                    if (cmdline) {
                        kfree(cmdline);
                    }
                }
            }
        }
    }

    seq_printf(m, "\n]\n}\n");
    return 0;
}

// Función que se ejecuta al abrir el archivo /proc
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

// Estructura que contiene las operaciones del archivo /proc
static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};

// Función de inicialización del módulo
static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "sysinfo_json modulo cargado\n");
    return 0;
}

// Función de limpieza del módulo
static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo_json modulo desinstalado\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);

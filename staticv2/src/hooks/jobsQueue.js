import { useCallback, useState } from 'react';

/**
 * Un Custom Hook para poder agregar jobs a una cola y ejecutarlos en orden.
 * @param {Number} ms - El tiempo en milisegundos que se debe esperar entre cada job.
 * @returns {Object} Un objeto con las funciones para agregar, deshacer y ejecutar jobs.
 */
const useJobsQueue = (ms) => {
    const [ jobsQueue, setJobsQueue ] = useState([]);
        
    /**
     * Agrega un job a la cola.
     * @param {Function} actionCb - La función que se debe ejecutar para aplicar el job.
     * @param {Function} undoCb - La función que se debe ejecutar para deshacer el job.
     * @param {Function} clearJob - La función que se debe ejecutar para limpiar el job (no lo corre).
     * @param {Object} actionInfo - La información del job.
     */
    const addJob = useCallback(({ actionCb, undoCb, clearCb, actionInfo }) => {
        const duration = ms + ((ms / 2) * jobsQueue.length);
        
        const timeoutId = setTimeout(() => {
            actionCb();
            setJobsQueue((jobsQueue) => jobsQueue.filter(job => job.id !== timeoutId));
        }, duration);

        const job = {
            id: timeoutId,
            undoJobFunc: undoCb,
            actionFunc: actionCb,
            clearFunc: clearCb,
            info: actionInfo,
            deleteIn: duration,
        };

        setJobsQueue((prevJobsQueue) => [...prevJobsQueue, job]);
        return job;
    }, [ ms, jobsQueue ]);
    
    /**
     * Deshace un job. Y lo elimina de la cola.
     * @param {Object} job - El job que se debe deshacer.
     */
    const undoJob = useCallback((job) => {
        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
        setJobsQueue((jobsQueue) => jobsQueue.filter(j => j.id !== job.id));
    }, []);

    /** 
     * Ejecuta un job. Y lo elimina de la cola.
     * @param {Object} job - El job que se debe ejecutar.
     */
    const executeJob = useCallback((job) => {
        clearTimeout(job.id);
        job.actionFunc();
        setJobsQueue((jobsQueue) => jobsQueue.filter(j => j.id !== job.id));
    }, []);

    /**
     * Limpia un job. Y lo elimina de la cola.
     * @param {Object} job - El job que se debe limpiar.
     */
    const clearJob = useCallback((job) => {
        clearTimeout(job.id);
        if(job.clearFunc) job.clearFunc();
        setJobsQueue((jobsQueue) => jobsQueue.filter(j => j.id !== job.id));
    }, []);

    /**
     * Deshace el último job en la cola.
     */
    const undoLastJob = useCallback(() => {
        if(jobsQueue.length === 0) return;
        const job = jobsQueue.shift();
        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
        setJobsQueue(jobsQueue);
    }, [ jobsQueue ]);

    /**
     * Borra todos los jobs en la cola.
     * No deshace los jobs.
     * Solo los borra.
     **/
    const clearAllJobs = useCallback(() => {
        if(jobsQueue.length === 0) return;

        jobsQueue.forEach(job => clearTimeout(job.id));
        jobsQueue.forEach(job => job.clearFunc && job.clearFunc());
        setJobsQueue([]);
    }, [ jobsQueue ]);

    /**
     * Deshace todos los jobs en la cola.
     */
    const undoAllJobs = useCallback(() => {
        if(jobsQueue.length === 0) return;

        // Borra todos los timeouts, luego deshace todos los jobs.
        jobsQueue.forEach(job => clearTimeout(job.id));
        jobsQueue.forEach(job => job.undoJobFunc && job.undoJobFunc());
        setJobsQueue([]);
    }, [ jobsQueue ]);

    /**
     * Ejecuta todos los jobs en la cola.
     */
    const executeAllJobs = useCallback(() => {
        if(jobsQueue.length === 0) return;
        
        // Borra todos los timeouts, luego deshace todos los jobs.
        jobsQueue.forEach(job => clearTimeout(job.id));
        jobsQueue.forEach(job => job.actionFunc());
        setJobsQueue([]);
    }, [ jobsQueue ]);

    return {
        addJob, 
        undoJob,
        executeJob,
        clearJob,
        undoLastJob,
        undoAllJobs,
        executeAllJobs,
        clearAllJobs,
        jobsQueue
    };
};

export default useJobsQueue;
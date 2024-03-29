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
     * @param {Function} actionCb - La función que se debe ejecutar.
     * @param {Function} undoCb - La función que se debe ejecutar para deshacer el job.
     * @param {Object} actionInfo - La información del job.
     */
    const addJob = useCallback((actionCb, undoCb, actionInfo) => {
        const timeoutId = setTimeout(() => {
            actionCb();
            setJobsQueue((jobsQueue) => jobsQueue.filter(job => job.id !== timeoutId));
        }, ms + (ms * jobsQueue.length));

        const job = {
            id: timeoutId,
            undoJobFunc: undoCb,
            action: actionCb,
            info: actionInfo
        };

        setJobsQueue((prevJobsQueue) => [...prevJobsQueue, job]);
        return job;
    }, [ ms, jobsQueue ]);
    
    /**
     * Deshace un job.
     * @param {Object} job - El job que se debe deshacer.
     */
    const undoJob = useCallback((job) => {
        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
    }, []);

    /**
     * Deshace el último job en la cola.
     */
    const undoLastJob = useCallback(() => {
        if(jobsQueue.length === 0) return;
        const job = jobsQueue.shift();
        setJobsQueue(jobsQueue);

        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
    }, [ jobsQueue ]);

    /**
     * Deshace todos los jobs en la cola.
     */
    const clearAllJobs = useCallback(() => {
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
        jobsQueue.forEach(job => job.action());
        setJobsQueue([]);
    }, [ jobsQueue ]);

    return {
        addJob, 
        undoLastJob,
        undoJob,
        clearAllJobs,
        executeAllJobs,
        jobsQueue
    };
};

export default useJobsQueue;
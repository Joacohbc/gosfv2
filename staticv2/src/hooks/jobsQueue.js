import { useCallback, useState } from 'react';
/**
 * Custom hook for managing a jobs queue with a specified delay.
 *
 * @param {number} ms - The delay in milliseconds for executing each job.
 * @returns {Object} - An object containing the addJob function, undoLastJob function, and jobsQueue array.
 */
const useJobsQueue = (ms) => {
    const [ jobsQueue, setJobsQueue ] = useState([]);
        
    // Add a job to the queue. (returns the new job)
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
    
    // Undo a specific job.
    const undoJob = useCallback((job) => {
        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
    }, []);

    // Undo the last job in the queue.
    const undoLastJob = useCallback(() => {
        if(jobsQueue.length === 0) return;
        const job = jobsQueue.shift();
        setJobsQueue(jobsQueue);

        clearTimeout(job.id);
        if(job.undoJobFunc) job.undoJobFunc();
    }, [ jobsQueue ]);

    // Undo all jobs in the queue.
    const clearAllJobs = useCallback(() => {
        if(jobsQueue.length === 0) return;

        // First clear all timeouts, then undo all jobs.
        jobsQueue.forEach(job => clearTimeout(job.id));
        jobsQueue.forEach(job => job.undoJobFunc && job.undoJobFunc());
        setJobsQueue([]);
    }, [ jobsQueue ]);

    const executeAllJobs = useCallback(() => {
        if(jobsQueue.length === 0) return;
        
        // First clear all timeouts, then execute all jobs.
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
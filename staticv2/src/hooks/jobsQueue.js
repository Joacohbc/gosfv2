import { useCallback,useState } from 'react';

/**
 * Custom hook for managing a jobs queue with a specified delay.
 *
 * @param {number} ms - The delay in milliseconds for executing each job.
 * @returns {Object} - An object containing the addJob function, undoLastJob function, and jobsQueue array.
 */
const useJobsQueue = (ms) => {
    const [ jobsQueue, setJobsQueue ] = useState([]);

    
    const addJob = useCallback((actionCb, undoCb) => {
        const timeoutId = setTimeout(() => {
            actionCb();
            setJobsQueue((jobsQueue) => jobsQueue.filter(job => job.timeoutId !== timeoutId));
        }, ms);

        setJobsQueue((prevJobsQueue) => {
            prevJobsQueue.unshift({
                timeoutId: timeoutId,
                undoJobFunc: undoCb
            });
            return prevJobsQueue;
        });
    }, [ ms ]);
    
    const undoLastJob = useCallback(() => {
        if(jobsQueue.length === 0) return;
        const job = jobsQueue.shift();
        clearTimeout(job.timeoutId);
        if(job.undoJobFunc) job.undoJobFunc();
    }, [ jobsQueue ]);

    return {
        addJob, 
        undoLastJob,
        jobsQueue
    };
};

export default useJobsQueue;
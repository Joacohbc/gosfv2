import { useCallback,useState } from 'react';

const useJobsQueue = (ms) => {
    const [ jobsQueue, setJobsQueue ] = useState([]);

    const addJob = useCallback((cb, undoCb) => {
        const timeoutId = setTimeout(() => {
            cb();
            setJobsQueue((jobsQueue) => jobsQueue.filter(job => job.timeoutId != timeoutId));
        }, ms);

        setJobsQueue((jobsQueue) => {
            jobsQueue.unshift({
                timeoutId: timeoutId,
                undoJobFunc: undoCb
            });
            return jobsQueue;
        });
    }, [ ms ]);
    
    const undoLastJob = useCallback(() => {
        if(jobsQueue.length == 0) return;
        const job = jobsQueue.shift();
        clearTimeout(job.timeoutId);
        job.undoJobFunc();
    }, [ jobsQueue ]);

    return {
        addJob, 
        undoLastJob,
        jobsQueue
    };
};

export default useJobsQueue;
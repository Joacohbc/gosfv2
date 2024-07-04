import { getCacheService } from '../services/cache';
import { useMemo } from 'react';

const ONE_SECOND = 1000;
const ONE_MINUTE = 60 * ONE_SECOND;

export const useCache = () => {
    return useMemo(() => ({
        cacheService: getCacheService(),
        ONE_SECOND,
        ONE_MINUTE
    }), []);
} 
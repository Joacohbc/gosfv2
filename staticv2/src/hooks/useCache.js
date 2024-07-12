import { getCacheService } from '../services/cache';
import { useMemo } from 'react';

export const useCache = () => {
    return useMemo(() => getCacheService(), []);
} 
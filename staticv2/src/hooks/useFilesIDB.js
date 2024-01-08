
import { useCallback, useEffect, useState } from 'react'

const useFilesIDB = () => {

    const [ db, setDb ] = useState(null)

    useEffect(() => {
        const request = window.indexedDB.open('files-db', 3)
        
        request.onerror = (event) => {
            console.log('Error opening indexedDB: ', event)
        }
        
        request.onsuccess = (event) => {
            setDb(event.target.result)
        }
        
        request.onupgradeneeded = (event) => {
            const db = event.target.result;

            const objectStore = db.createObjectStore('files', { keyPath: 'id', autoIncrement: false })
            objectStore.createIndex('id', 'id', { unique: true })
            objectStore.createIndex('filename', 'filename', { unique: false })
            objectStore.createIndex('blob', 'blob', { unique: false })
        }

    }, [ ]);

    const addFileToLocal = useCallback((id, filename, blob) => {
        return new Promise((resolve, reject) => {
            if(!db) reject('No indexedDB')

            const transaction = db.transaction('files', 'readwrite')
            const objectStore = transaction.objectStore('files')

            const request = objectStore.add({ id, filename, blob })

            request.onsuccess = (event) => {
                resolve(event.target.result)
            }

            request.onerror = (event) => {
                reject(event);
            }
        });
    }, [ db ]);

    const getFileFromLocal = useCallback((fileId) => {
        return new Promise((resolve, reject) => {
            if(!db) reject('No indexedDB')

            const transaction = db.transaction('files', 'readonly')
            const objectStore = transaction.objectStore('files')

            const request = objectStore.get(fileId)

            request.onsuccess = (event) => {
                resolve(event.target.result)
            }

            request.onerror = (event) => {
                reject(event)
            }
        });
    }, [ db ]);

    const deleteFileFromLocal = useCallback((fileId) => {
        return new Promise((resolve, reject) => {
            if(!db) reject('No indexedDB')

            const transaction = db.transaction('files', 'readwrite')
            const objectStore = transaction.objectStore('files')

            const request = objectStore.delete(fileId)

            request.onsuccess = (event) => {
                resolve(event.target.result)
            }

            request.onerror = (event) => {
                reject(event)
            }
        });
    }, [ db ]);

    const deleteAllFilesFromLocal = useCallback(() => {
        return new Promise((resolve, reject) => {
            if(!db) reject('No indexedDB')

            const transaction = db.transaction('files', 'readwrite')
            const objectStore = transaction.objectStore('files')

            const request = objectStore.clear()

            request.onsuccess = (event) => {
                resolve(event.target.result)
            }

            request.onerror = (event) => {
                reject(event)
            }
        });
    }, [ db ]);
    
    return {
        addFileToLocal,
        getFileFromLocal,
        deleteFileFromLocal,
        deleteAllFilesFromLocal
    }
}   

export default useFilesIDB
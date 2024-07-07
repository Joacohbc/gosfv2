import './SearchBar.css';
import { useCallback, useState } from 'react';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import PropTypes from "prop-types";
import ToolTip from './ToolTip';

const sortByName = (order) => {
    return (data) => data.sort((a, b) => {
        if(order == 'asc') {
            if(a.filename < b.filename) {
                return -1;
            }

            if(a.filename > b.filename) {
                return 1;
            }

            return 0;
        }

        if(order == 'desc') {
            if(a.filename > b.filename) {
                return -1;
            }

            if(a.filename < b.filename) {
                return 1;
            }

            return 0;
        }
    });
}

const filterByText = (text) => {
    return (data) => data.filter(file => file.filename.toLowerCase().includes(text.toLowerCase()) || file.id == text);
}

const sortByDate = (order) => {
    return (data) => data.sort((a, b) => {
        const dateA = new Date(a.createdAt);
        const dateB = new Date(b.createdAt);

        if(order == 'asc') {
            if(dateA < dateB) {
                return -1;
            }

            if(dateA > dateB) {
                return 1;
            }

            return 0;
        }

        if(order == 'desc') {
            if(dateA > dateB) {
                return -1;
            }

            if(dateA < dateB) {
                return 1;
            }
        }

        return 0;
    });
}

const SearchBar = ({ createFileLoader, setFileLoader  }) => {
    const [ sortName, setSortName ] = useState('asc');
    const [ disableSort, setDisableSort ] = useState(false);
    const [ sortDate, setSortDate ] = useState('asc');

    const handleSearch = handleKeyUpWithTimeout(async (e) => {
        setDisableSort(e.target.value != '');
        const loadInfo = await createFileLoader(filterByText(e.target.value));
        setFileLoader(loadInfo);
        loadInfo();
    }, 500);

    const handleSortByName = useCallback(async () => {
        const newSortName = sortName == 'asc' ? 'desc' : 'asc';
        setSortName(newSortName);

        const loadInfo = await createFileLoader(sortByName(newSortName));
        setFileLoader(loadInfo);
        loadInfo();
    }, [ sortName, setSortName, setFileLoader, createFileLoader ]);

    const handleSortByDate = useCallback(async () => {
        const newSortDate = sortDate == 'asc' ? 'desc' : 'asc';
        setSortDate(newSortDate);

        const loadInfo = await createFileLoader(sortByDate(newSortDate));
        setFileLoader(loadInfo);
        loadInfo();
    }, [ sortDate, setSortDate, setFileLoader, createFileLoader ])

    return (
        <div className="d-flex justify-content-center align-items-center mb-4 gap-1">
            <input type="text" placeholder="Enter Search" className='search-input' onKeyUp={handleSearch}/>
            <div className={`sort-icons ${disableSort && 'sort-icons-disable'}`}>
                { sortName == 'asc' && 
                <ToolTip toolTipMessage='Sort by name A to Z' placement='top'>
                    <i className='bi bi-sort-alpha-down' onClick={handleSortByName}/>
                </ToolTip> }
                
                { sortName == 'desc' &&
                <ToolTip toolTipMessage='Sort by name Z to A' placement='top'>
                    <i className='bi bi-sort-alpha-up' onClick={handleSortByName}/>
                </ToolTip> }
            </div>

            <div className={`sort-icons ${disableSort && 'sort-icons-disable'}`}>
                { sortDate == 'asc' && <ToolTip toolTipMessage='Sort by date oldest to newest' placement='top'>
                    <i className='bi bi-sort-numeric-down' onClick={handleSortByDate}/>
                </ToolTip> }

                { sortDate == 'desc' && <ToolTip toolTipMessage='Sort by date newest to oldest' placement='top'>
                    <i className='bi bi-sort-numeric-up' onClick={handleSortByDate}/>
                </ToolTip> }
            </div>
        </div>
    );
}

SearchBar.propTypes = {
    createFileLoader: PropTypes.func.isRequired,
    setFileLoader: PropTypes.func.isRequired
};

export default SearchBar;
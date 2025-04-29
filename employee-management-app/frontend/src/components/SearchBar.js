import React, { useState } from 'react';

export default function SearchBar({value, setValue, searchType, setSearchType, onSearch}){   
    function handleChange(e){
        setValue(e.target.value);
    }
    
    return (
        <div className="search-container">
            <section className="search-bar">
                <i className="fa-solid fa-magnifying-glass search-icon"></i>
                <input 
                    className="search-input" 
                    placeholder="Search by employee..."
                    value={value}
                    onChange={handleChange}
                    onKeyDown={onSearch}
                />
            </section>
            
            <button 
                className="search-button"
                onMouseDown={e => e.preventDefault()}
                onClick={onSearch}
            >
                Search
            </button>
        </div>
    );
}
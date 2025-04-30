import React, { useState } from 'react';

export default function SearchBar({value, setValue, refreshEmployeeData}) {   
    function handleChange(e){
        setValue(e.target.value);
    }
    
    const handleKeyUp = e => {
        if (e.key === 'Enter') {
          refreshEmployeeData();
        }
    };


    return (
        <div className="search-container">
            <section className="search-bar">
                <i className="fa-solid fa-magnifying-glass search-icon"></i>
                <input 
                    className="search-input" 
                    placeholder="Search by employee..."
                    value={value}
                    onChange={handleChange}
                    onKeyUp={handleKeyUp}
                />
            </section>
            
            <button 
                className="search-button"
                onMouseDown={e => e.preventDefault()}
                onClick={refreshEmployeeData}
            >
                Search
            </button>
        </div>
    );
}
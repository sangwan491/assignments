import React from 'react';

export default function SearchBar({value, setValue}){
    function handleChange(e){
        setValue(e.target.value);
    }
    return (
        <section className="search-bar">
            <i className="fa-solid fa-magnifying-glass search-icon"></i>
            <input 
                className="search-input" 
                placeholder="Search by employee name/id..." 
                value={value} 
                onChange={handleChange}
            />
        </section>
    );
}
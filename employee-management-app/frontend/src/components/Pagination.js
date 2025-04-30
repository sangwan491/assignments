import React from 'react';

const Pagination = ({ currentPage, setCurrentPage, pagesArray, totalPages }) => {
  const goToPreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const goToNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
    }
  };

  return <>
    { totalPages != 0 &&
      (
        <section className='pagination'>
          <button 
            className='pagination-button' 
            onClick={goToPreviousPage}
            disabled={currentPage === 1}
          >
            <i className="fa-solid fa-arrow-left"></i>
          </button>
          <section className='pagination-numbers'>
            {pagesArray.map(page => (
              <span 
                key={page} 
                onClick={() => setCurrentPage(page)} 
                className={`pagination-number ${currentPage === page ? 'active' : ''}`}
              >
                {page}
              </span>
            ))}
          </section>
          <button 
            className='pagination-button'
            onClick={goToNextPage}
            disabled={currentPage === totalPages}
          >
            <i className="fa-solid fa-arrow-right"></i>
          </button>
        </section>
      )
    } 
  </>
};

export default Pagination;

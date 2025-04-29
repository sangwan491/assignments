import React, { useState, useEffect, useMemo, useCallback } from 'react';
import SearchBar from './components/SearchBar';
import EmployeeForm from './components/EmployeeForm';
import EmployeeList from './components/EmployeeList';
import Pagination from './components/Pagination';
import './App.css';

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const API_BASE_URL = 'http://localhost:8080/api';

function App() {
  const [isAdding, setAdding] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  // const [searchType, setSearchType] = useState("employee");
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [searchTerm, setSearchTerm] = useState("");
  const [employeePages, setEmployeesPages] = useState({});
  const [fetchKey, setFetchKey] = useState(0);

  const fetchUrl = useMemo(() => {
    return `${API_BASE_URL}/employees/search?term=${searchTerm}&page=${currentPage}`;
  }, [currentPage, searchTerm]);

  const employees = useMemo(() => {
    return employeePages[currentPage] || [];
  }, [currentPage, employeePages]);

  const pagesArray = useMemo(() => {
    if (!totalPages) {
      return [];
    }

    let start, end;
    if (totalPages <= 5) {
      start = 1;
      end = totalPages;
    } else {
      if (currentPage <= 3) {
        start = 1;
        end = 5;
      } else if (currentPage >= totalPages - 2) {
        start = totalPages - 4;
        end = totalPages;
      } else {
        start = currentPage - 2;
        end = currentPage + 2;
      }
    }
    return Array.from({ length: end - start + 1 }, (_, index) => start + index);
  }, [currentPage, totalPages]);

  const handleSearch = useCallback(async () => {
    setIsLoading(true);

      try {
        const response = await fetch(fetchUrl);
        const jsonResponse = await response.json();
  
        if (!response.ok) {
          throw new Error(jsonResponse.error || `Search failed with status: ${response.status}`);
        }
  
        const results = jsonResponse.employees;

        const pages = jsonResponse.pages;
  
        setTotalPages(pages);
        setEmployeesPages(prevPages => ({
          ...prevPages,
          [currentPage]: results 
        }));
      } catch (err) {
        console.error('Search error:', err);
        toast.error(`Search failed: ${err.message}`);
      } finally {
        setIsLoading(false);
      }
  }, [employeePages, fetchUrl]);

  const refreshEmployeeData = useCallback(() => {
    setEmployeesPages({});
    setFetchKey(prev => prev + 1);
  }, []);
  
  useEffect(() => {
    if (employeePages[currentPage]) {
      return;
    }

    handleSearch();
  }, [currentPage, fetchKey]);

  const onSave = async (newEmployee) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newEmployee),
      });
      
      const jsonResponse = await response.json();
      
      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }

      refreshEmployeeData();
      setCurrentPage(1);
      
      toast.success(jsonResponse.message);

      setAdding(false);  
    } catch (err) {
      console.log(err);
      toast.error(`Failed to add employee: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };
  
  const onCancel = () => {
    setAdding(false);
  };


  const onUpdate = async (updatedEmployee) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees/${updatedEmployee.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedEmployee),
      });
      
      const jsonResponse = await response.json();

      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }

      // Update the employee in the cache. Loop through each page.
      setEmployeesPages((prevPages) => {
        const newPages = { ...prevPages };
        Object.keys(newPages).forEach(page => {
          newPages[page] = newPages[page].map(emp =>
            emp.id === updatedEmployee.id ? updatedEmployee : emp
          );
        });
        return newPages;
      });

      toast.success(jsonResponse.message);
    } catch (err) {
      console.log(`err: ${err}`);
      toast.error(`Failed to update employee: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };


  const onDelete = async (id) => {
    toast.promise(
      deleteEmployee(id),
      {
        pending: 'Deleting employee...',
        success: 'Employee deleted successfully',
        error: {
          render({data}) {
            return data.message;
          }
        }
      }
    );
  };
  
  const deleteEmployee = async (id) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${API_BASE_URL}/employees/${id}`, {
        method: 'DELETE',
      });
      
      const jsonResponse = await response.json();

      if (!response.ok) {
        throw new Error(jsonResponse.error);
      }
      
      refreshEmployeeData();
      setCurrentPage(1);

      return jsonResponse;
    } catch (err) {
      console.log(err);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };


  return (
    <>      
      <header>
        Employee Management
      </header>
      
      <main className="app-content">
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />

        {isAdding ? (
          <EmployeeForm 
            employees={employees}
            onSave={onSave} 
            onCancel={onCancel}
            isLoading={isLoading}
          />
        ) : (
          <button 
            className='add-new-btn' 
            onClick={() => setAdding(true)}
            disabled={isLoading}
          >
            Add new Employee
          </button>
        )}
        
        <SearchBar 
          value={searchTerm} 
          setValue={setSearchTerm} 
          setCurrentPage={setCurrentPage}
          onSearch={refreshEmployeeData}
        />
        
        {isLoading? (
          <span className="loading">Loading employees...</span>
        ) : (
          <EmployeeList
            employees={employees}
            onUpdate={onUpdate}
            onDelete={onDelete}
            isLoading={isLoading}
          />
        )}
      </main>

      <footer>
        <Pagination currentPage={currentPage} setCurrentPage={setCurrentPage} pagesArray={pagesArray} totalPages={totalPages}></Pagination>  
      </footer>
    </>
  );
}

export default App;
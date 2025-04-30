import React, { useState, useEffect, useMemo } from 'react';
import SearchBar from './components/SearchBar';
import EmployeeForm from './components/EmployeeForm';
import EmployeeList from './components/EmployeeList';
import Pagination from './components/Pagination';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';

import { searchEmployees, createEmployee, updateEmployee, deleteEmployee } from './utils/api';
import { calculatePagesArray } from './utils/pagination';
import { validateEmployee } from './utils/validation';    

function App() {
  const [isAdding, setAdding] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [searchTerm, setSearchTerm] = useState("");
  const [employeePages, setEmployeesPages] = useState({});
  const [fetchKey, setFetchKey] = useState(0);

  const employees = useMemo(() => {
    return employeePages[currentPage] || [];
  }, [currentPage, employeePages]);

  const pagesArray = useMemo(() => 
    calculatePagesArray(currentPage, totalPages),
  [currentPage, totalPages]);

  const handleSearch = async () => {
    setIsLoading(true);
    try {
      const jsonResponse = await searchEmployees(searchTerm, currentPage);
      setTotalPages(jsonResponse.pages);
      setEmployeesPages(prevPages => ({
        ...prevPages,
        [currentPage]: jsonResponse.employees
      }));
    } catch (err) {
      console.error('Search error:', err);
      toast.error(`Search failed: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const refreshEmployeeData = () => {
    setEmployeesPages({});
    setFetchKey(prev => prev + 1);
  };
  
  useEffect(() => {
    if (employeePages[currentPage]) {
      return;
    }
    handleSearch();
  }, [currentPage, fetchKey]);

  const onSave = async (newEmployee) => {
    const validationErrors = validateEmployee(newEmployee);
    if (validationErrors.length > 0) {
      validationErrors.forEach(error => toast.error(error));
      return;
    }

    setIsLoading(true);
    try {
      const response = await createEmployee(newEmployee);
      refreshEmployeeData();
      setCurrentPage(1);
      toast.success(response.message);
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
    const validationErrors = validateEmployee(updatedEmployee);
    if (validationErrors.length > 0) {
      validationErrors.forEach(error => toast.error(error));
      return;
    }

    setIsLoading(true);
    try {
      const response = await updateEmployee(updatedEmployee);
      setEmployeesPages((prevPages) => {
        const newPages = { ...prevPages };
        Object.keys(newPages).forEach(page => {
          newPages[page] = newPages[page].map(emp =>
            emp.id === updatedEmployee.id ? updatedEmployee : emp
          );
        });
        return newPages;
      });
      toast.success(response.message);
    } catch (err) {
      console.log(`err: ${err}`);
      toast.error(`Failed to update employee: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const onDelete = async (id) => {
    toast.promise(
      deleteEmployee(id).then(() => {
        refreshEmployeeData();
        setCurrentPage(1);
      }),
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
          refreshEmployeeData={refreshEmployeeData}
        />
        
        {isLoading ? (
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
        <Pagination 
          currentPage={currentPage} 
          setCurrentPage={setCurrentPage} 
          pagesArray={pagesArray} 
          totalPages={totalPages}
        />
      </footer>
    </>
  );
}

export default App;
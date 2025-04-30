import { API_BASE_URL } from './constants';
import { toast } from 'react-toastify';

export const searchEmployees = async (searchTerm, currentPage) => {
    const response = await fetch(`${API_BASE_URL}/employees/search?term=${searchTerm}&page=${currentPage}`);
    const jsonResponse = await response.json();

    if (!response.ok) {
        throw new Error(jsonResponse.error || `Search failed with status: ${response.status}`);
    }

    return jsonResponse;
};

export const createEmployee = async (newEmployee) => {
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

    return jsonResponse;
};

export const updateEmployee = async (updatedEmployee) => {
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

    return jsonResponse;
};

export const deleteEmployee = async (id) => {
    const response = await fetch(`${API_BASE_URL}/employees/${id}`, {
        method: 'DELETE',
    });
    
    const jsonResponse = await response.json();

    if (!response.ok) {
        throw new Error(jsonResponse.error);
    }

    return jsonResponse;
};
export const validateEmployee = (employee) => {
    const errors = [];
    
    if (!employee.name || employee.name.trim() === '') {
      errors.push('Name is required');
    }
    
    if (!employee.email || employee.email.trim() === '') {
      errors.push('Email is required');
    } else if (!/\S+@\S+\.\S+/.test(employee.email)) {
      errors.push('Email format is invalid');
    }
    
    if (!employee.department || employee.department.trim() === '') {
      errors.push('Department is required');
    }
    
    return errors;
};
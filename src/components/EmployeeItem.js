import React, { useEffect, useState } from "react";

const EmployeeItem = ({ employee, onUpdate, onDelete }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedEmployee, setEditedEmployee] = useState({ ...employee });
  const [isDeleting, setIsDeleting] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setEditedEmployee(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSave = () => {
    onUpdate(editedEmployee);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditedEmployee({ ...employee });
    setIsEditing(false);
  };

  useEffect(() => {
    setEditedEmployee({ ...employee });
  }, [employee]);

  return (
    <div className="employee-item">
      <div className="employee-header">
        <input
          type="text"
          name="name"
          value={editedEmployee.name}
          onChange={handleChange}
          disabled={!isEditing}
      
        />
        <div className="employee-actions">
          {isEditing ? (
            <>
              <button className="btn-save" onClick={handleSave}>
                Save
              </button>
              <button className="btn-cancel" onClick={handleCancel}>
                Cancel
              </button>
            </>
          ) : isDeleting ? (
            <>
              <button className="btn-confirm" onClick={() => onDelete(employee.id)}>
                <i class="fa-solid fa-check"></i>
              </button>
              <button className="btn-cancel" onClick={() => setIsDeleting(false)}>
                <i class="fa-solid fa-xmark"></i>
              </button>
            </>
          ) : (
            <>
              <button className="btn-edit" onClick={() => setIsEditing(true)}>
                <i className="fa-solid fa-pen-to-square"></i>
              </button>
              <button className="btn-delete" onClick={() => setIsDeleting(true)}>
                <i className="fa-solid fa-trash"></i>
              </button>
            </>
          )
        }
        </div>
      </div>
      <hr />
      <div className="employee-details">
        <div className="detail-item">
          <i className="fa-solid fa-envelope"></i>
          <span className="label">Email:</span>
          <input 
            type="email"
            name="email"
            value={editedEmployee.email}
            onChange={handleChange}
            disabled={!isEditing}
          />
        </div>
        <div className="detail-item">
          <i className="fa-solid fa-phone"></i>
          <span className="label">Phone:</span>
          <input 
            type="text"
            name="phone"
            value={editedEmployee.phone}
            onChange={handleChange}
            disabled={!isEditing}
          />
        </div>
        <div className="detail-item">
          <i className="fa-solid fa-user"></i>
          <span className="label">Department:</span>
          <input 
            type="text"
            name="department"
            value={editedEmployee.department}
            onChange={handleChange}
            disabled={!isEditing}
          />
        </div>
      </div>
    </div>
  );
};

export default EmployeeItem;
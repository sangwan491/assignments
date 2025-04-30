export const calculatePagesArray = (currentPage, totalPages) => {
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
};
type PaginationProps = {
  page: () => number
  handlePrevPage: () => void
  handleNextPage: () => void
}

export default function Pagination({ page, handlePrevPage, handleNextPage }: PaginationProps) {
  return (
    <div class="bg-white rounded-lg shadow p-4 flex items-center justify-between">
      <button
        class="px-3 py-1 border rounded text-sm disabled:opacity-50"
        disabled={page() <= 1}
        onClick={handlePrevPage}
      >
        Previous
      </button>
      <span class="text-sm">Page {page()}</span>
      <button
        class="px-3 py-1 border rounded text-sm"
        onClick={handleNextPage}
      >
        Next
      </button>
    </div>
  )
}
import { createSignal, createMemo, JSX, createEffect } from "solid-js";
import { ColumnDef } from "@tanstack/solid-table";
import { Checkbox } from '~/components/ui/checkbox';

export function useTableSelection<T extends { id?: string }>() {
  const [selectedIds, setSelectedIds] = createSignal<Set<string>>(new Set());

  const toggleSelection = (id: string) => {
    setSelectedIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const toggleAll = (data: T[]) => {
    setSelectedIds((prev) => {
      const allIds = data.map((item) => item.id).filter((id): id is string => id !== undefined);
      const allSelected = allIds.every((id) => prev.has(id));

      if (allSelected) {
        // Deselect all
        const newSet = new Set(prev);
        allIds.forEach((id) => newSet.delete(id));
        return newSet;
      } else {
        // Select all
        const newSet = new Set(prev);
        allIds.forEach((id) => newSet.add(id));
        return newSet;
      }
    });
  };

  const clearSelection = () => {
    setSelectedIds(new Set<string>());
  };

  const isSelected = (id: string) => selectedIds().has(id);

  const selectedCount = () => selectedIds().size;

  const getSelectedItems = <T extends { id?: string }>(data: T[]): T[] => {
    return data.filter((item) => item.id !== undefined && selectedIds().has(item.id));
  };

  const createSelectionColumn = <T extends { id?: string }>(
    data: () => T[]
  ): ColumnDef<T> => ({
    id: "select",
    header: () => {
      // Create reactive getters for the checkbox state
      const isAllSelected = createMemo(() => {
        const currentData = data();
        const allIds = currentData.map((item) => item.id).filter((id): id is string => id !== undefined);
        const selected = selectedIds();
        return allIds.length > 0 && allIds.every((id) => selected.has(id));
      });
      
      const isSomeSelected = createMemo(() => {
        const currentData = data();
        const allIds = currentData.map((item) => item.id).filter((id): id is string => id !== undefined);
        const selected = selectedIds();
        return allIds.some((id) => selected.has(id)) && !isAllSelected();
      });

      return (
        <Checkbox
          checked={isAllSelected()}
          indeterminate={isSomeSelected()}
          onChange={() => toggleAll(data())}
          aria-label="Select all"
        />
      );
    },
    cell: (info) => {
      const item = info.row.original;
      const id = item.id;
      if (id === undefined) return null;
      return (
        <Checkbox
          checked={isSelected(id)}
          onChange={() => toggleSelection(id)}
          aria-label={`Select row`}
          onClick={(e: MouseEvent) => e.stopPropagation()}
        />
      );
    },
    size: 40,
    enableSorting: false,
  });

  return {
    selectedIds,
    toggleSelection,
    toggleAll,
    clearSelection,
    isSelected,
    selectedCount,
    getSelectedItems,
    createSelectionColumn,
  };
}

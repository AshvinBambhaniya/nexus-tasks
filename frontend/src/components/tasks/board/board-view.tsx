"use client";

import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  useSensors,
  useSensor,
  PointerSensor,
  DragStartEvent,
  DragOverEvent,
  closestCorners
} from "@dnd-kit/core";
import { useState, useEffect } from "react";
import { Task, TaskStatus } from "@/types";
import { BoardColumn } from "./board-column";
import { BoardCard } from "./board-card";
import { createPortal } from "react-dom";
import { arrayMove } from "@dnd-kit/sortable";

interface BoardViewProps {
  tasks: Task[];
  onTaskMove: (taskId: number, newStatus: TaskStatus) => void;
  onTaskClick?: (task: Task) => void;
}

const COLUMNS = [
  { id: TaskStatus.BACKLOG, title: "Backlog", color: "bg-gray-400" },
  { id: TaskStatus.TODO, title: "To Do", color: "bg-blue-500" },
  { id: TaskStatus.IN_PROGRESS, title: "In Progress", color: "bg-yellow-500" },
  { id: TaskStatus.DONE, title: "Done", color: "bg-green-500" },
];

export function BoardView({ tasks: initialTasks, onTaskMove, onTaskClick }: BoardViewProps) {
  const [tasks, setTasks] = useState(initialTasks);
  const [activeTask, setActiveTask] = useState<Task | null>(null);
  
  useEffect(() => {
      setTasks(initialTasks);
  }, [initialTasks]);
  
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const handleDragStart = (event: DragStartEvent) => {
    setActiveTask(event.active.data.current?.task);
  };

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over } = event;
    if (!over) return;

    const activeId = active.id;
    const overId = over.id;

    // If moving over a different column (container), update status optimistically
    // Note: DND Kit Sortable handles reordering, but moving between containers needs state update if we want "live" preview
    // For simplicity, we stick to DragEnd for status changes, but DragOver helps visual placement
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTask(null);

    if (!over) return;

    const activeTask = active.data.current?.task;
    const overId = over.id;

    // Check if dropped on a column (status change)
    if (Object.values(TaskStatus).includes(overId as TaskStatus)) {
        if (activeTask.status !== overId) {
            onTaskMove(activeTask.id, overId as TaskStatus);
        }
        return;
    }

    // Check if dropped on another task (reorder or status change if diff column)
    const overTask = tasks.find(t => t.id === overId);
    if (overTask) {
        if (activeTask.status !== overTask.status) {
             onTaskMove(activeTask.id, overTask.status);
        } else {
            // Reorder logic (local only since backend doesn't support it yet)
            // We can update local state to show reordering visually
            const oldIndex = tasks.findIndex((t) => t.id === active.id);
            const newIndex = tasks.findIndex((t) => t.id === over.id);
            setTasks((items) => arrayMove(items, oldIndex, newIndex));
        }
    }
  };

  return (
    <DndContext 
        sensors={sensors} 
        collisionDetection={closestCorners}
        onDragStart={handleDragStart} 
        onDragOver={handleDragOver}
        onDragEnd={handleDragEnd}
    >
      <div className="flex h-full gap-4 overflow-x-auto pb-4">
        {COLUMNS.map((col) => (
          <BoardColumn
            key={col.id}
            id={col.id}
            title={col.title}
            color={col.color}
            tasks={tasks.filter((t) => t.status === col.id)}
            onTaskClick={onTaskClick}
          />
        ))}
      </div>
      
      {typeof document !== "undefined" && createPortal(
        <DragOverlay>
          {activeTask ? <BoardCard task={activeTask} /> : null}
        </DragOverlay>,
        document.body
      )}
    </DndContext>
  );
}

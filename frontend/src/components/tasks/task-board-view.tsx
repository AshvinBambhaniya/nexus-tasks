"use client";

import {
  DndContext,
  DragEndEvent,
} from "@dnd-kit/core";
import { Task, TaskStatus } from "@/types";
import { BoardColumn } from "./board-column";

const COLUMNS = [
  { id: TaskStatus.TODO, title: "To Do" },
  { id: TaskStatus.IN_PROGRESS, title: "In Progress" },
  { id: TaskStatus.DONE, title: "Done" },
];

interface TaskBoardViewProps {
  tasks: Task[];
  onTaskMove: (taskId: number, newStatus: TaskStatus) => void;
  onTaskClick?: (task: Task) => void;
}

export function TaskBoardView({ tasks, onTaskMove, onTaskClick }: TaskBoardViewProps) {
  
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (!over) return;

    const taskId = active.id as number;
    const newStatus = over.id as TaskStatus;

    if (active.data.current?.status !== newStatus) {
        onTaskMove(taskId, newStatus);
    }
  };

  return (
    <DndContext onDragEnd={handleDragEnd}>
      <div className="flex h-full gap-6 overflow-x-auto pb-4">
        {COLUMNS.map((col) => (
          <BoardColumn
            key={col.id}
            id={col.id}
            title={col.title}
            tasks={tasks.filter((t) => t.status === col.id)}
            onTaskClick={onTaskClick}
          />
        ))}
      </div>
    </DndContext>
  );
}
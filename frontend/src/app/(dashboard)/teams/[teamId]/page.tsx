"use client";

import { use, useState } from "react";
import { useTeams, useTeamMembers } from "@/hooks/use-teams";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { UserPlus } from "lucide-react";

export default function TeamPage({ params }: { params: Promise<{ teamId: string }> }) {
  const resolvedParams = use(params);
  const teamId = parseInt(resolvedParams.teamId);
  
  const { teams } = useTeams();
  const { members, isLoading, addMember } = useTeamMembers(teamId);
  const team = teams.find(t => t.id === teamId);

  const [inviteEmail, setInviteEmail] = useState("");
  const [isInviting, setIsInviting] = useState(false);

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteEmail) return;
    
    setIsInviting(true);
    try {
      await addMember(inviteEmail);
      setInviteEmail("");
    } catch (error) {
      alert("Failed to add member");
    } finally {
      setIsInviting(false);
    }
  };

  return (
    <div className="flex flex-col h-full space-y-6">
       <div className="space-y-1">
           <h1 className="text-2xl font-bold tracking-tight text-gray-900">
             {team ? team.name : "Team"}
           </h1>
           {team?.description && (
             <p className="text-sm text-gray-500">{team.description}</p>
           )}
        </div>

        <div className="grid gap-6 md:grid-cols-2">
            <Card className="p-6">
                <h3 className="text-lg font-medium mb-4">Members ({members.length})</h3>
                <div className="space-y-4">
                    {isLoading ? (
                        <p>Loading...</p>
                    ) : (
                        members.map((member: any) => (
                            <div key={member.user_id} className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                    <Avatar className="h-8 w-8" fallback={member.email[0].toUpperCase()} />
                                    <div>
                                        <p className="text-sm font-medium">{member.email}</p>
                                        <p className="text-xs text-gray-500">{member.role}</p>
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </Card>

            <Card className="p-6 h-fit">
                <h3 className="text-lg font-medium mb-4">Add Member</h3>
                <form onSubmit={handleInvite} className="flex gap-2">
                    <Input 
                        placeholder="user@example.com" 
                        value={inviteEmail}
                        onChange={(e) => setInviteEmail(e.target.value)}
                        required
                        type="email"
                    />
                    <Button type="submit" disabled={isInviting}>
                        <UserPlus className="h-4 w-4 mr-2" />
                        Add
                    </Button>
                </form>
            </Card>
        </div>
    </div>
  );
}

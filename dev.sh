#!/bin/bash
# Copyright (C) 2025 Michael Graff
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, version 3.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.


set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Cleanup function
cleanup() {
    echo -e "\n${BLUE}Shutting down...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start backend
echo -e "${GREEN}Starting Go backend...${NC}"
export JWT_SECRET="dev-secret-key-change-in-production"
export DB_USER="$USER"
export DB_NAME="meshmgr"
export DB_SSLMODE="disable"
./bin/meshmgr &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

# Start frontend
echo -e "${GREEN}Starting SvelteKit frontend...${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!

echo -e "\n${GREEN}Services running:${NC}"
echo -e "  Backend:  http://localhost:8080"
echo -e "  Frontend: http://localhost:5173"
echo -e "\n${BLUE}Press Ctrl+C to stop${NC}\n"

# Wait for either process to exit
wait

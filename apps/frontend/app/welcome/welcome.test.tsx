import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Welcome } from "~/welcome/welcome";

describe("Welcome screen", () => {
    it("should show the welcome screen", async () => {
        render(<Welcome />);

        expect(screen.getByText(/what's next/i)).toBeInTheDocument();
    });
});
import React from 'react';

interface HeaderProps {
  title: string;
  subtitle: string;
}

export function Header({ title, subtitle }: HeaderProps) {
  return (
    <header className="text-center py-0 mb-[67px]">
      <h1 className="text-[44px] font-bold text-white tracking-wide mb-2">
        <span className="bg-gradient-to-r from-[#ffd36b] via-[#fbbf24] to-[#f59e0b] bg-clip-text text-transparent">
          {title}
        </span>
      </h1>
      <p className="text-white/70 text-sm tracking-widest font-normal">
        {subtitle}
      </p>
    </header>
  );
}
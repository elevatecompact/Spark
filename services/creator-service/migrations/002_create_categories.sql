CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    icon_url TEXT DEFAULT '',
    color VARCHAR(7) DEFAULT '#6366F1',
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    sort_order INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_sort_order ON categories(sort_order ASC);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

INSERT INTO categories (name, slug, description, icon_url, color, sort_order, active) VALUES
('Music', 'music', 'Musicians, singers, and bands', '', '#EF4444', 1, true),
('Gaming', 'gaming', 'Gamers and streamers', '', '#8B5CF6', 2, true),
('Technology', 'technology', 'Tech reviews, coding, and gadgets', '', '#3B82F6', 3, true),
('Education', 'education', 'Tutorials and educational content', '', '#10B981', 4, true),
('Entertainment', 'entertainment', 'Comedy, skits, and variety shows', '', '#F59E0B', 5, true),
('Sports', 'sports', 'Athletes and sports commentary', '', '#06B6D4', 6, true),
('News & Politics', 'news-politics', 'News coverage and political analysis', '', '#DC2626', 7, true),
('Lifestyle', 'lifestyle', 'Daily life, vlogs, and wellness', '', '#EC4899', 8, true),
('Art & Design', 'art-design', 'Visual arts, design, and illustration', '', '#F97316', 9, true),
('Fashion & Beauty', 'fashion-beauty', 'Style, makeup, and fashion tips', '', '#D946EF', 10, true),
('Food & Cooking', 'food-cooking', 'Cooking shows, recipes, and food reviews', '', '#EAB308', 11, true),
('Travel', 'travel', 'Travel vlogs and destination guides', '', '#14B8A6', 12, true),
('Fitness & Health', 'fitness-health', 'Workouts, health tips, and nutrition', '', '#22C55E', 13, true),
('Science', 'science', 'Scientific explanations and experiments', '', '#6366F1', 14, true),
('Business & Finance', 'business-finance', 'Entrepreneurship, investing, and finance', '', '#059669', 15, true),
('Film & Animation', 'film-animation', 'Movies, animation, and filmmaking', '', '#7C3AED', 16, true),
('DIY & Crafts', 'diy-crafts', 'Do-it-yourself projects and crafting', '', '#F43F5E', 17, true),
('Pets & Animals', 'pets-animals', 'Pet care, animal facts, and cute content', '', '#84CC16', 18, true),
('Comedy', 'comedy', 'Stand-up, sketches, and funny content', '', '#F97316', 19, true),
('Dance', 'dance', 'Dance performances and tutorials', '', '#E11D48', 20, true),
('Photography', 'photography', 'Photography tips, gear, and showcases', '', '#0EA5E9', 21, true),
('Writing & Literature', 'writing-literature', 'Writing tips, book reviews, and poetry', '', '#8B5CF6', 22, true),
('Podcasts', 'podcasts', 'Podcast episodes and audio content', '', '#6B7280', 23, true),
('Spirituality', 'spirituality', 'Meditation, religion, and spiritual growth', '', '#A855F7', 24, true),
('Automotive', 'automotive', 'Cars, motorcycles, and vehicle reviews', '', '#18181B', 25, true),
('Outdoors', 'outdoors', 'Hiking, camping, and outdoor adventures', '', '#15803D', 26, true),
('Home & Garden', 'home-garden', 'Home improvement, gardening, and decor', '', '#65A30D', 27, true),
('Parenting', 'parenting', 'Parenting tips, family vlogs, and child care', '', '#FB923C', 28, true),
('Anime & Manga', 'anime-manga', 'Anime reviews, manga, and otaku culture', '', '#EC4899', 29, true),
('Other', 'other', 'Content that does not fit other categories', '', '#9CA3AF', 999, true);
